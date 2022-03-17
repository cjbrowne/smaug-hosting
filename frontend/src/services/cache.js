let inMemoryCache = {};

let Get = (key) => {
    if (inMemoryCache[key]) {
        return inMemoryCache[key];
    }

    if (localStorage.getItem(key) !== null) {
        return localStorage.getItem(key);
    }

    return null;
};

let Put = (key, value) => {
    if(!value) {
        localStorage.removeItem(key);
        delete inMemoryCache[key];
    } else {
        inMemoryCache[key] = value;
        localStorage.setItem(key, value);
    }
};

export {
    Get,
    Put
}