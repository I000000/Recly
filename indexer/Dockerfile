FROM python:3.12-slim
RUN pip config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple && \
    pip install meilisearch pandas pyarrow
COPY index_data.py /app/index_data.py
CMD ["python", "/app/index_data.py"]