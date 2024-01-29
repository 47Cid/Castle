import logging
from datetime import datetime
from flask import Flask, jsonify, request
app = Flask(__name__)


# Set up logging
logging.basicConfig(filename='logs/http.log', level=logging.INFO)

# Basic SQL injection detection


@app.route('/verify')
def verify():
    logging.info('[%s] Received request: %s %s\nHeaders: %s\nBody: %s',
                 datetime.now().isoformat(), request.method, request.url, request.headers, request.get_data())

    body = request.get_data().decode('utf-8')

    sql_keywords = ['SELECT', 'DROP', ';', '--', '\'']
    if any(keyword in body.upper() for keyword in sql_keywords):
        logging.warning('[%s] Possible SQL injection attack: %s',
                        datetime.now().isoformat(), body)
        return jsonify(isValid='False', time=datetime.now().isoformat())

    return jsonify(isValid='True', time=datetime.now().isoformat())


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=3031)
