![Alt text](/assets/WAF.png)


###  An asynchronous, batteries-not-included WAF framework that can be easily customized and extended.

## How it works:
1) The reverse proxy first intercepts the client's connection.
2) It then inspects the destination URL and routes the request to a specific "pod". This pod is responsible for conducting various validation checks on the request.
3) If the request passes all validation checks and is deemed non-malicious, it is then relayed to the server.
4) Throughout this process, both the proxy and the individual pods generate and maintain logs for tracking and auditing purposes.