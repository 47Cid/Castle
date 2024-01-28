from flask import Flask, jsonify
app = Flask(__name__)
from datetime import datetime


@app.route('/verify')
def verify():
    return jsonify(isValid='True', time=datetime.now().isoformat())

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=3031)