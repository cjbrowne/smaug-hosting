import url from 'url';
import {getToken} from "./auth";

let baseUrl = process.env.REACT_APP_BASE_CONTAINER_SERVICE_URL;

let create = ({name, software, tier}) => {
    return new Promise((resolve, reject) => {
        fetch(url.resolve(baseUrl, '/containers/'), {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${getToken()}`
            },
            body: JSON.stringify({
                name,
                software,
                tier
            })
        }).then((response) => {
            response.json().then(response.ok ? resolve : reject, reject);
        }, reject)
    });
};

let getAll = () => {
    return new Promise((resolve, reject) => {
        fetch(url.resolve(baseUrl, '/containers/'), {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${getToken()}`
            }
        }).then((response) => {
            response.json().then(response.ok ? resolve : reject, reject);
        }, reject)
    });
};

let destroy = (id) => {
    return new Promise((resolve, reject) => {
        fetch(url.resolve(baseUrl, `/containers/${id}/`), {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${getToken()}`
            }
        }).then((response) => {
            response.json().then(response.ok ? resolve : reject, reject);
        }, reject);
    });
};

let stop = (id) => {
    return new Promise((resolve, reject) => {
        fetch(url.resolve(baseUrl, `/containers/${id}/stop/`), {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${getToken()}`
            }
        }).then((response) => {
            response.json().then(response.ok ? resolve : reject, reject);
        }, reject)
    })
};

let start = (id) => {
    return new Promise((resolve, reject) => {
        fetch(url.resolve(baseUrl, `/containers/${id}/start/`), {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${getToken()}`
            }
        }).then((response) => {
            response.json().then(response.ok ? resolve : reject, reject);
        }, reject);
    });
};

export {
    create,
    destroy,
    getAll,
    stop,
    start
}