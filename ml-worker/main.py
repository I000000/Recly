import os, json, signal, sys, time
import numpy as np, pandas as pd, h5py
import psycopg2, redis, pika
from sklearn.metrics.pairwise import cosine_similarity

RABBITMQ_URL = os.getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
REDIS_URL = os.getenv("REDIS_URL", "redis://localhost:6379/0")
QUEUE_NAME = os.getenv("QUEUE_NAME", "recommendation_tasks")
EMBEDDINGS_PATH = os.getenv("EMBEDDINGS_PATH", "/data/embeddings.h5")
BOOK_PARQUET = os.getenv("BOOK_PARQUET", "/data/goodreads.parquet")
MOVIE_PARQUET = os.getenv("MOVIE_PARQUET", "/data/tmdb.parquet")

DATABASE_URL = os.getenv("DATABASE_URL", "postgres://recly:recly_pass@postgres:5432/recly_db?sslmode=disable")

book_text = movie_text = book_genre = movie_genre = book_image = movie_image = None
book_ids = movie_ids = None
redis_client = None

def load_ids(path, id_col):
    return pd.read_parquet(path, columns=[id_col])[id_col].astype(str).values

def load_embeddings():
    global book_text, movie_text, book_genre, movie_genre, book_image, movie_image, book_ids, movie_ids
    with h5py.File(EMBEDDINGS_PATH, 'r') as f:
        book_text  = f['/books/desc'][:].astype('float32')
        book_genre = f['/books/genre'][:].astype('float32')
        book_image = f['/books/image'][:].astype('float32')
        movie_text  = f['/movies/desc'][:].astype('float32')
        movie_genre = f['/movies/genre'][:].astype('float32')
        movie_image = f['/movies/image'][:].astype('float32')
    book_ids  = load_ids(BOOK_PARQUET, 'book_id')
    movie_ids = load_ids(MOVIE_PARQUET, 'movie_id')
    print(f"Loaded embeddings for {len(book_ids)} books and {len(movie_ids)} movies.")

def normalize(vec):
    norm = np.linalg.norm(vec)
    return vec / norm if norm > 1e-8 else vec

def recommend(selected_indices, weights, direction):
    w_g, w_t, w_i = weights.get('genre',0.3), weights.get('text',0.4), weights.get('image',0.3)
    src = ('book' if direction.startswith('book') else 'movie')
    tgt = ('movie' if direction.endswith('movie') else 'book')
    src_text, src_genre, src_img, src_ids = (book_text, book_genre, book_image, book_ids) if src == 'book' else (movie_text, movie_genre, movie_image, movie_ids)
    tgt_text, tgt_genre, tgt_img, tgt_ids = (movie_text, movie_genre, movie_image, movie_ids) if tgt == 'movie' else (book_text, book_genre, book_image, book_ids)

    user_text = normalize(sum(src_text[i] for i in selected_indices))
    user_genre = normalize(sum(src_genre[i] for i in selected_indices))
    user_img = normalize(sum(src_img[i] for i in selected_indices))

    sim_text  = cosine_similarity(user_text.reshape(1,-1), tgt_text).flatten()
    sim_genre = cosine_similarity(user_genre.reshape(1,-1), tgt_genre).flatten()
    sim_img   = cosine_similarity(user_img.reshape(1,-1), tgt_img).flatten()
    combined  = w_g*sim_genre + w_t*sim_text + w_i*sim_img
    top = np.argsort(combined)[::-1][:10]
    return [tgt_ids[i] for i in top]

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
        id_to_idx = {bid: i for i, bid in enumerate(book_ids if direction.startswith('book') else movie_ids)}
        selected_indices = [id_to_idx[sid] for sid in data['selected_ids'] if sid in id_to_idx]
        if not selected_indices:
            result = {"status": "error", "error": "No valid selected items"}
        else:
            result = {"status": "done", "movies": recommend(selected_indices, weights, direction)}
        redis_client.set(f"rec:{task_id}", json.dumps(result), ex=1800)
        print(f"Task {task_id} completed. Recommendations are ready.")
        if result.get("status") == "done":
            update_history(task_id, json.dumps(result["movies"]))
    except Exception as e:
        print(f"Error: {e}")
        if task_id:
            redis_client.set(f"rec:{task_id}", json.dumps({"status":"error","error":str(e)}), ex=1800)
    finally:
        ch.basic_ack(delivery_tag=method.delivery_tag)

def start_consumer():
    """Подключается к RabbitMQ и запускает потребление с переподключением при ошибках."""
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