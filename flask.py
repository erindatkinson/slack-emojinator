"""dev flask app"""
from flask import Flask, request

app = Flask(__name__)

@app.route("/test", methods=["POST"])
def test():
    """checker for post vars"""
    print(request.path)
    print(request.method)
    print(request.headers)
    print(request.form)
    print(request.files)
    return "ok"
