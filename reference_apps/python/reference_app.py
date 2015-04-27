import json
from flask import Flask, session, redirect, request

app = Flask(__name__)
app.secret = 'insecureappsecret'
app.debug = True

CAS_LOGIN_ADDR = 'http://localhost:3000/login?service=localhost:3001/validateCASLogin'

@app.route('/', methods=['GET'])
def index():
    return "Hello %s! <br/> <a href=\"/login\">Login?</a>" % (session['user']['name'] if 'user' in session else 'stranger')

@app.route('/login', methods=['GET'])
def login():
    return redirect(CAS_LOGIN_ADDR)

@app.route('/validateCASLogin', methods=['GET'])
def cas_validate():
    user = request.get_json(force=True)
    session['user'] = user
    return redirect()

if __name__ == '__main__':
    app.run(port=3001)
