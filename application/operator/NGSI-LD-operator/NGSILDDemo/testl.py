from flask import Flask, jsonify, abort, request, make_response
import requests 
import json
app = Flask(__name__, static_url_path = "")

@app.route('/notifyContext', methods = ['POST'])
def hello():
	print "=============notify============="
	if not request.json:
		abort(400)
	print(request.json)
	return "Hello"

if __name__ == "__main__":
	app.run(host='0.0.0.0', port=8888)    
