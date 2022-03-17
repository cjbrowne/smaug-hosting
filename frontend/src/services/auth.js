import { Get, Put } from './cache';
import url from 'url';

let baseUrl = process.env.REACT_APP_BASE_AUTH_URL;

let hasAuth = () => {
    // todo: handle expired tokens proactively
    return Get('token') != null;
};

let setToken = (token) => {
    Put('token', token.Token);
    Put('refresh', token.Refresh);
    Put('expires', token.Expires);
};

let getToken = () => {
    return Get('token');
};

let logout = () => {
    setToken({});
};

let register = (email, password) => {
    return new Promise((resolve, reject) => {
        fetch(url.resolve(baseUrl, '/user/'), {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                email,
                password
            })
        }).then((response) => {
            response.json().then(response.ok ? resolve : reject, reject);
        }, reject);
    });
};

let verify = (verificationToken) => {
    return new Promise((resolve, reject) => {
        fetch(url.resolve(baseUrl, `/verify/${verificationToken}/`)).then((response) => {
            response.json().then(response.ok ? resolve : reject, reject);
        }, reject);
    });
};

let login = (email, password) => {
    return new Promise((resolve, reject) => {
        fetch(url.resolve(baseUrl, '/token/'), {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                email,
                password
            })
        }).then((response) => {
            response.json().then(response.ok ? (token) => {
                setToken(token);
                resolve(token);
            } : reject, reject);
        }, reject);
    });
};

export {
    hasAuth,
    login,
    logout,
    register,
    getToken,
    verify
}