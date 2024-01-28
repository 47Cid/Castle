import logging
from datetime import datetime
from flask import Flask, jsonify, request
app = Flask(__name__)


# Set up logging
logging.basicConfig(filename='logs/http.log', level=logging.INFO)


@app.route('/verify')
def verify():
    logging.info('[%s] Received request: %s %s\nHeaders: %s\nBody: %s',
                 datetime.now().isoformat(), request.method, request.url, request.headers, request.get_data())

    return jsonify(isValid='True', time=datetime.now().isoformat())


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=3031)
