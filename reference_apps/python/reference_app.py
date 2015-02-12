from flask import Flask

app = Flask(__name__)

@app.route('/', methods=['GET'])
def index():
    return 'Hello World!'



if __name__ == '__main__':
    app.run(port=3001)
