import json, urllib, requests
from flask import Flask, session, redirect, request, url_for

app = Flask(__name__)
app.secret_key = 'insecureappsecretkey'
app.debug = True

# NOTE - these endpoints SHOULD be using HTTPS (for communication with the CAS server)
CAS_ADDR = 'http://localhost:3000/'
SERVICE_URL = 'http://localhost:3001/validateCASLogin'
CAS_LOGIN_ADDR = "".join([CAS_ADDR, "/login?service=", urllib.quote_plus(SERVICE_URL)])
CAS_CHECK_ADDR_TEMPLATE = "".join([CAS_ADDR, "/validate?", "service=", SERVICE_URL, "&ticket=%s"])

@app.route('/', methods=['GET'])
def index():
    if 'userEmail' in session and 'userAttributes' in session:
        return "<h2>Welcome user with email %s!</h2> <p> Attributes: <code>%s</code> </p>" % (session['userEmail'], session['userAttributes'])
    else:
        return "Hello %s! <br/> <a href=\"/login\">Login?</a>" % (session['userEmail'] if 'userEmail' in session else 'stranger')

@app.route('/login', methods=['GET'])
def login():
    return redirect(CAS_LOGIN_ADDR)

@app.route('/logout', methods=['GET'])
def logout():
    if 'userEmail' in session and 'userAttributes' in session:
        del session['userEmail']
        del session['userAttributes']
    return 'Successfully logged out'

@app.route('/validateCASLogin', methods=['GET'])
def cas_validate():
    ticket = request.args['ticket']

    # Lookup ticket with CAS Server
    lookup_addr = CAS_CHECK_ADDR_TEMPLATE % ticket
    print("Lookup addr %s" % lookup_addr)
    cas_resp = requests.get(lookup_addr).json()
    print("Resp:", cas_resp)

    # Error handling
    if cas_resp['status'] == 'error':
        return "Oh No! An error ocurred:<br/> <strong>%s</strong>" % cas_resp['message']
    else:
        session['userEmail'] = cas_resp['userEmail']
        session['userAttributes'] = cas_resp['userAttributes']
        
    return redirect(url_for('index'))

if __name__ == '__main__':
    app.run(port=3001)
