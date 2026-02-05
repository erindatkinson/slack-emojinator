FROM python:3.14-alpine
RUN pip install pipenv
WORKDIR /emojinator
VOLUME /emojinator/export
VOLUME /emojinator/import
COPY Pipfile .
COPY Pipfile.lock .
RUN pipenv install
COPY __init__.py .
COPY main.py .
COPY packages/ ./packages/
ENTRYPOINT [ "pipenv", "run", "python"]