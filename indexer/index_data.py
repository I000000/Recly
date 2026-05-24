import meilisearch, pandas as pd, numpy as np, time, os, json, sys, requests

MEILI_URL = os.getenv("MEILI_URL", "http://meilisearch:7700")
MEILI_KEY = os.getenv("MEILI_KEY", "aSecretMasterKey")
BOOK_PQ = os.getenv("BOOK_PARQUET", "/data/goodreads_2k.parquet")
MOVIE_PQ = os.getenv("MOVIE_PARQUET", "/data/tmdb_2k.parquet")
BATCH_SIZE = 1000

client = meilisearch.Client(MEILI_URL, MEILI_KEY)

print("Waiting for Meilisearch...", flush=True)
for _ in range(60):
    try:
        client.health()
        print("Meilisearch is ready.", flush=True)
        break
    except:
        time.sleep(2)

def clean_for_json(obj):
    if obj is None: return None
    if isinstance(obj, float):
        if np.isnan(obj) or np.isinf(obj): return None
        return obj
    if isinstance(obj, (int, bool, str)): return obj
    if isinstance(obj, np.ndarray): return clean_for_json(obj.tolist())
    if isinstance(obj, list): return [clean_for_json(item) for item in obj]
    if isinstance(obj, dict): return {k: clean_for_json(v) for k, v in obj.items()}
    return str(obj)

def fix_documents(df):
    df = df.where(pd.notnull(df), None)
    for col in df.select_dtypes(include=['datetime64', 'datetimetz']).columns:
        df[col] = df[col].apply(lambda x: x.isoformat() if x else None)
    records = df.to_dict(orient='records')
    for rec in records:
        for k in rec:
            rec[k] = clean_for_json(rec[k])
    return records

def add_in_batches(documents, label):
    total = len(documents)
    for start in range(0, total, BATCH_SIZE):
        batch = documents[start:start+BATCH_SIZE]
        try:
            payload = json.dumps(batch)
        except Exception as e:
            print(f"  {label}: batch {start}-{start+len(batch)} JSON error: {e}", flush=True)
            for i, doc in enumerate(batch):
                try: json.dumps(doc)
                except: print(f"    Problem doc at index {start+i}: {doc.get('id', 'no id')}", flush=True)
            continue

        url = f"{MEILI_URL}/indexes/items/documents?primaryKey=id"
        headers = {"Content-Type": "application/json", "Authorization": f"Bearer {MEILI_KEY}"}
        try:
            resp = requests.post(url, data=payload, headers=headers)
            if resp.status_code not in (200, 202):
                try:
                    error_detail = resp.json()
                    print(f"  {label}: batch {start}-{start+len(batch)} HTTP {resp.status_code}: {error_detail.get('message', resp.text)}", flush=True)
                except:
                    print(f"  {label}: batch {start}-{start+len(batch)} HTTP {resp.status_code}: {resp.text}", flush=True)
                continue
            task = resp.json()
            while True:
                st = requests.get(f"{MEILI_URL}/tasks/{task['taskUid']}", headers=headers).json()
                if st['status'] in ('succeeded', 'failed'): break
                time.sleep(0.5)
            if st['status'] == 'failed':
                print(f"  {label}: batch {start}-{start+len(batch)} FAILED: {st.get('error', {})}", flush=True)
            else:
                print(f"  {label}: {min(start+BATCH_SIZE, total)}/{total}", flush=True)
        except Exception as e:
            print(f"  {label}: batch {start}-{start+len(batch)} HTTP error: {e}", flush=True)

# --- КНИГИ ---
print("Loading books...", flush=True)
books = pd.read_parquet(BOOK_PQ,
                         columns=['book_id','title','authors','image_url','genres','average_rating','description','publication_year'])
books['type'] = 'book'
books['id'] = books['book_id'].astype(str)
if 'authors' in books.columns:
    books['authors'] = books['authors'].apply(lambda a: ', '.join([d.get('name', str(d)) for d in a]) if isinstance(a, list) else str(a) if a else None)
documents = fix_documents(books)
add_in_batches(documents, "Books")
print(f"All books processed ({len(documents)} total).", flush=True)

# --- ФИЛЬМЫ ---
print("Loading movies...", flush=True)
movies = pd.read_parquet(MOVIE_PQ,
                         columns=['movie_id','title','poster_full_url','genres','vote_average','director','cast','overview','release_date','runtime'])
movies['type'] = 'movie'
movies['id'] = movies['movie_id'].astype(str)
movies['image_url'] = movies['poster_full_url']
if 'cast' in movies.columns:
    movies['cast'] = movies['cast'].apply(lambda c: ', '.join(c) if isinstance(c, list) else str(c) if c else None)
if 'release_date' in movies.columns:
    movies['release_date'] = movies['release_date'].apply(lambda x: x.isoformat() if hasattr(x, 'isoformat') else str(x) if x else None)
documents = fix_documents(movies)
add_in_batches(documents, "Movies")
print(f"All movies processed ({len(documents)} total).", flush=True)

# --- НАСТРОЙКИ ---
print("Updating index settings...", flush=True)
client.index('items').update_searchable_attributes(['title'])
client.index('items').update_filterable_attributes(['id', 'type'])
client.index('items').update_displayed_attributes(['*'])
print("Indexing completed.", flush=True)
