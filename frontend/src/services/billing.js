import _ from 'lodash';
import {getToken} from "./auth";
import url from 'url';

let baseBillingUrl = process.env.REACT_APP_BASE_BILLING_URL;
// todo: env var for this
let billingWsUrl = process.env.REACT_APP_BILLING_WEBSOCKET_URL;

let backoff = 1000;

let initWs = () => {
    let ws = new WebSocket(billingWsUrl);
    ws.onopen = () => {
        ws.send(JSON.stringify({
            subject: "handshake",
            body: {
                "token": getToken()
            }
        }));
    };

    ws.onclose = () => {
        setTimeout(initWs, backoff *= 2);
    };
    return ws;
};

let ws = initWs();

let balanceSubscribers = {};

let updateBalanceSubscribers = (balance) => {
    _.each(balanceSubscribers, (sub) => {
        sub(balance);
    })
};

ws.onmessage = (message) => {
    let payload = JSON.parse(message.data);
    switch(payload.subject) {
        case "balance": {
            updateBalanceSubscribers(payload.body.balance);
            break;
        }
        default: {
            break;
        }
    }
};

let subscribeToBalance = (callback) => {
    // todo: build a more "safe" random id generator
    let id = Math.random();

    balanceSubscribers[id] = callback;

    return () => {
        delete balanceSubscribers[id];
    }
};

let topup = (amount) => {
    console.log(process.env);
    return new Promise((resolve, reject) => {
        fetch(url.resolve(baseBillingUrl, `/topup/${amount}/`), {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${getToken()}`
            },
        }).then((response) => {
            response.json().then(response.ok ? resolve : reject, reject);
        }, reject);
    });
};

export {
    subscribeToBalance,
    topup
}