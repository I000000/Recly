import os, json, signal, sys, time, random
import numpy as np, pandas as pd, h5py
import psycopg2, redis, pika

RABBITMQ_URL = os.getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
REDIS_URL = os.getenv("REDIS_URL", "redis://localhost:6379/0")
QUEUE_NAME = os.getenv("QUEUE_NAME", "recommendation_tasks")
EMBEDDINGS_PATH = os.getenv("EMBEDDINGS_PATH", "/data/embeddings.h5")
BOOK_PARQUET = os.getenv("BOOK_PARQUET", "/data/goodreads.parquet")
MOVIE_PARQUET = os.getenv("MOVIE_PARQUET", "/data/tmdb.parquet")
DATABASE_URL = os.getenv("DATABASE_URL", "postgres://recly:recly_pass@postgres:5432/recly_db?sslmode=disable")

book_text = movie_text = book_genre = movie_genre = book_image = movie_image = None
book_ids = movie_ids = None
book_id_to_idx = {}
movie_id_to_idx = {}
redis_client = None

def load_ids(path, id_col):
    return pd.read_parquet(path, columns=[id_col])[id_col].astype(str).values

def normalize_matrix(m):
    """L2-нормализация матрицы построчно (in-place)"""
    norms = np.linalg.norm(m, axis=1, keepdims=True)
    norms[norms == 0] = 1
    return m / norms

def load_embeddings():
    print(f"Loading embeddings. Please wait.")
    global book_text, movie_text, book_genre, movie_genre, book_image, movie_image, book_ids, movie_ids
    global book_id_to_idx, movie_id_to_idx

    with h5py.File(EMBEDDINGS_PATH, 'r') as f:
        book_text  = f['/books/desc'][:].astype('float32')
        book_genre = f['/books/genre'][:].astype('float32')
        book_image = f['/books/image'][:].astype('float32')
        movie_text  = f['/movies/desc'][:].astype('float32')
        movie_genre = f['/movies/genre'][:].astype('float32')
        movie_image = f['/movies/image'][:].astype('float32')

    print("Normalizing embeddings...")
    book_text   = normalize_matrix(book_text)
    book_genre  = normalize_matrix(book_genre)
    book_image  = normalize_matrix(book_image)
    movie_text  = normalize_matrix(movie_text)
    movie_genre = normalize_matrix(movie_genre)
    movie_image = normalize_matrix(movie_image)

    book_ids  = load_ids(BOOK_PARQUET, 'book_id')
    movie_ids = load_ids(MOVIE_PARQUET, 'movie_id')
    book_id_to_idx = {bid: i for i, bid in enumerate(book_ids)}
    movie_id_to_idx = {mid: i for i, mid in enumerate(movie_ids)}

    print(f"Loaded embeddings: books {book_text.shape}, movies {movie_text.shape}.")
    print(f"Book IDs: {len(book_ids)}, Movie IDs: {len(movie_ids)}.")

def normalize_vector(vec):
    norm = np.linalg.norm(vec)
    return vec / norm if norm > 1e-8 else vec

def recommend(selected_ids, weights, direction, exclude_ids=None, selected_weights=None):
    start_time = time.time()

    w_g = weights.get('genre', 0.3)
    w_t = weights.get('text', 0.4)
    w_i = weights.get('image', 0.3)

    user_text = np.zeros(book_text.shape[1], dtype=np.float32)
    user_genre = np.zeros(book_genre.shape[1], dtype=np.float32)
    user_image = np.zeros(book_image.shape[1], dtype=np.float32)
    weight_sum = 0.0

    for sid in selected_ids:
        w = selected_weights.get(sid, 1.0) if selected_weights else 1.0
        if sid.startswith('book_'):
            idx = book_id_to_idx.get(sid[5:])
            if idx is not None:
                user_text += w * book_text[idx]
                user_genre += w * book_genre[idx]
                user_image += w * book_image[idx]
                weight_sum += w
        elif sid.startswith('movie_'):
            idx = movie_id_to_idx.get(sid[6:])
            if idx is not None:
                user_text += w * movie_text[idx]
                user_genre += w * movie_genre[idx]
                user_image += w * movie_image[idx]
                weight_sum += w

    if weight_sum > 0:
        user_text /= weight_sum
        user_genre /= weight_sum
        user_image /= weight_sum

    user_text = normalize_vector(user_text)
    user_genre = normalize_vector(user_genre)
    user_image = normalize_vector(user_image)

    if direction.endswith('movie'):
        tgt_text, tgt_genre, tgt_img, tgt_ids = movie_text, movie_genre, movie_image, movie_ids
        id_to_idx = movie_id_to_idx
    else:
        tgt_text, tgt_genre, tgt_img, tgt_ids = book_text, book_genre, book_image, book_ids
        id_to_idx = book_id_to_idx

    sim_text  = np.dot(user_text, tgt_text.T)
    sim_genre = np.dot(user_genre, tgt_genre.T)
    sim_img   = np.dot(user_image, tgt_img.T)

    combined = w_g * sim_genre + w_t * sim_text + w_i * sim_img

    if exclude_ids:
        for eid in exclude_ids:
            idx = id_to_idx.get(eid)
            if idx is not None:
                combined[idx] = -1e9

    top_indices = np.argsort(combined)[::-1][:30]
    result_ids = [tgt_ids[i] for i in top_indices]

    elapsed = (time.time() - start_time) * 1000  # мс
    print(f"Recommendation inference took {elapsed:.2f} ms")
    return result_ids

def update_history(task_id, movies_json):
    try:
        conn = psycopg2.connect(DATABASE_URL)
        cur = conn.cursor()
        cur.execute(
            "UPDATE user_recommendation_history SET result = %s::jsonb WHERE task_id = %s",
            (movies_json, task_id)
        )
        conn.commit()
        cur.close()
        conn.close()
        print(f"Task {task_id} saved to history.")
    except Exception as e:
        print(f"Failed to save task {task_id} to history: {e}.")

def on_message(ch, method, properties, body):
    task_id = None
    try:
        data = json.loads(body)
        task_id = data['task_id']
        direction = data.get('direction', 'book_to_movie')
        weights = data.get('weights', {})
        selected_ids = data['selected_ids']
        selected_weights = data.get('selected_weights', {})
        user_id = data.get('user_id')
        contextual = data.get('contextual', False)

        exclude_set = set()
        for sid in selected_ids:
            if sid.startswith('book_'):
                exclude_set.add(sid[5:])
            elif sid.startswith('movie_'):
                exclude_set.add(sid[6:])
        for eid in data.get('exclude_ids', []):
            exclude_set.add(str(eid))

        if user_id and user_id.strip() and not contextual:
            try:
                conn = psycopg2.connect(DATABASE_URL)
                cur = conn.cursor()
                cur.execute(
                    """SELECT result FROM user_recommendation_history
                       WHERE user_id = %s
                         AND result IS NOT NULL
                         AND jsonb_typeof(result) = 'array'
                         AND created_at > NOW() - INTERVAL '7 days'
                       ORDER BY created_at DESC""",
                    (user_id,)
                )
                added = 0
                for row in cur.fetchall():
                    try:
                        raw = row[0]
                        if isinstance(raw, list):
                            recent_ids = raw
                        else:
                            recent_ids = json.loads(raw)
                        for rid in recent_ids:
                            if rid and str(rid) not in exclude_set:
                                exclude_set.add(str(rid))
                                added += 1
                    except Exception as parse_err:
                        print(f"Warning: could not parse recent recommendation: {parse_err}")
                cur.close()
                conn.close()
                if added > 0:
                    print(f"Added {added} recent recommendations to exclude set")
            except Exception as e:
                print(f"Warning: could not load recent recommendations: {e}")

        base_w = {
            'text': weights.get('text', 0.4),
            'genre': weights.get('genre', 0.3),
            'image': weights.get('image', 0.3)
        }
        if contextual:
            noisy = base_w
        else:
            noise = {k: random.uniform(-0.05, 0.05) for k in base_w}
            noisy = {k: max(0.0, base_w[k] + noise[k]) for k in base_w}
            total = sum(noisy.values())
            if total > 0:
                noisy = {k: v / total for k, v in noisy.items()}
            else:
                noisy = base_w

        if not selected_ids:
            result = {"status": "error", "error": "No selected items"}
        else:
            result = {"status": "done", "movies": recommend(selected_ids, noisy, direction, exclude_ids=exclude_set, selected_weights=selected_weights)}

        redis_client.set(f"rec:{task_id}", json.dumps(result), ex=1800)
        print(f"Task {task_id} completed. Recommendations are ready.")

        if result.get("status") == "done" and not contextual:
            update_history(task_id, json.dumps(result["movies"]))
    except Exception as e:
        print(f"Error: {e}", flush=True)
        if task_id:
            redis_client.set(f"rec:{task_id}", json.dumps({"status":"error","error":str(e)}), ex=1800)
    finally:
        ch.basic_ack(delivery_tag=method.delivery_tag)

def start_consumer():
    while True:
        try:
            params = pika.URLParameters(RABBITMQ_URL)
            connection = pika.BlockingConnection(params)
            channel = connection.channel()
            channel.queue_declare(queue=QUEUE_NAME, durable=True)
            channel.basic_qos(prefetch_count=1)
            channel.basic_consume(queue=QUEUE_NAME, on_message_callback=on_message)
            print("Worker started. Waiting for messages...")
            channel.start_consuming()
        except (pika.exceptions.AMQPConnectionError, pika.exceptions.StreamLostError) as e:
            print(f"Connection lost: {e}. Reconnecting in 5 seconds...")
            time.sleep(5)
        except KeyboardInterrupt:
            raise

def shutdown(signum, frame):
    sys.exit(0)

def main():
    global redis_client
    redis_client = redis.Redis.from_url(REDIS_URL)
    print("Connected to Redis.")
    load_embeddings()
    signal.signal(signal.SIGINT, shutdown)
    signal.signal(signal.SIGTERM, shutdown)
    start_consumer()

if __name__ == "__main__":
    main()