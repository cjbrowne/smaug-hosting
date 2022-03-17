import _ from 'lodash';

let formatCurrency = (locale, number) => {
    return (new Intl.NumberFormat(locale.name, { style: 'currency', currency: locale.currency})).format(number);
};

let getCurrencyForCountry = (country) => {
    switch(country.toLowerCase()) {
        default:
        case 'gb':
            return 'GBP';
        case 'se':
            return 'GBP';
    }
};

let makeLocale = (language, country) => {
    // we don't support greek here, boiz
    return {
        name:`${language}-${country.toUpperCase()}`,
        currency: getCurrencyForCountry(country)
    };
};

let translations = {
    gamer: {
        generic: {
            name: "Name",
            software: "Game",
            tier: "Awesomeness",
            actions: "Press me",
            "ip-port": "Connection Information"
        },
        account: {
            topup: "Add moneyz",
            balance: "Fundz",
        },
        site: {
            title: "Smaug Hosting"
        },
        form: {
            auth: {
                title: "Let me in",
                email: {
                    label: "Sockpuppet"
                },
                password: {
                    label: "Secret Hacker Password"
                }
            }
        },
        whelps: {
            down: "unavailable"
        },
        containers: {
            name: "Whelp",
            plural: "Whelps",
            "create-success": "Whelp hatched!",
            minecraft: {
                description: "Get ready to grief 12 year olds",
                title: "Minecraft"
            },
            "create-dialog": {
                title: "Create Whelp",
                form: {
                    name: {
                        label: "Nick",
                        placeholder: "Whelpy McWhelpface"
                    },
                    tier: {
                        label: "Awesomeness Level"
                    }
                }
            },
            tiers: {
                0: {
                    name: "N00b",
                    help: "10 players max"
                },
                1: {
                    name: "pr0",
                    help: "25 player limit"
                },
                2: {
                    name: "Ultra pr0",
                    help: "30 player limit"
                },
                3: {
                    name: "Uberg33k",
                    help: "50 player limit"
                }
            }
        },
        cta: {
            login: "Go!",
            register: "Make new sockpuppet",
            "create-container": "Spin up a whelp!",
            create: "Make it",
            cancel: "Forget it",
            delete: "Delet this",
            stop: "Staahhp",
        },
        i18n: {
            "select-language": {
                label: "Language"
            }
        },
        errors: {
            login: {
                failed: "Nope"
            },
            label: "Fuck-up",
            "fetch-whelps": "I lost the list with your whelps on it"
        },
        nav: {
            home: "Home",
            "whelp-list": "Whelps",
            settings: "Config",
            logout: "GTFO"
        }
    },
    en: {
        settings: {
            title: "Settings",
            save: "Save"
        },
        generic: {
            name: "Name",
            software: "Software/Game",
            tier: "Tier",
            actions: "Actions",
            status: "Status",
            "ip-port": "Connection Information"
        },
        account: {
            topup: "Top up",
            balance: "Balance",
        },
        site: {
            title: "Smaug Hosting"
        },
        form: {
            auth: {
                title: "Login",
                email: {
                    label: "Email Address"
                },
                password: {
                    label: "Password"
                }
            }
        },
        "top-up": {
            title: "Top Up Credit",
            "amount-label": "Amount",
        },
        whelps: {
            down: "unavailable"
        },
        containers: {
            name: "Whelp",
            plural: "Whelps",
            "create-success": "Whelp hatched!",
            minecraft: {
                description: "Spin up a Minecraft server for you and your friends!",
                title: "Minecraft"
            },
            "create-dialog": {
                title: "Create Whelp",
                form: {
                    name: {
                        label: "Name",
                        placeholder: "Name your whelp!"
                    },
                    tier: {
                        label: "Tier"
                    }
                }
            },
            tiers: {
                0: {
                    name: "Fledgling",
                    help: "10 player limit"
                },
                1: {
                    name: "Lesser Dragon",
                    help: "25 player limit"
                },
                2: {
                    name: "Dragon",
                    help: "30 player limit"
                },
                3: {
                    name: "Greater Dragon",
                    help: "50 player limit"
                }
            }
        },
        cta: {
            login: "Login",
            register: "Register",
            "create-container": "Create Whelp",
            create: "Create",
            cancel: "Cancel",
            delete: "Delete",
            stop: "Stop",
            confirm: "Sure?",
            start: "Start",
            topup: "Top up",
        },
        i18n: {
            "select-language": {
                label: "Language"
            }
        },
        errors: {
            login: {
                failed: "Login Failed"
            },
            label: "Error",
            "fetch-whelps": "Could not fetch whelps",
            "stop-whelp": "Could not stop whelp",
            "delete-whelp": "Could not delete whelp"
        },
        nav: {
            home: "Home",
            "whelp-list": "My Whelps",
            settings: "Settings",
            logout: "Logout"
        }
    },
    se: {
        form: {
            auth: {
                title: "Logga in",
                email: {
                    label: "Email Address"
                },
                password: {
                    label: "Pass"
                }
            }
        },
        cta: {
            login: "Logga in",
            register: "Registrera"
        },
        i18n: {
            "select-language": {
                label: "Språk"
            }
        },
        errors: {
            login: {
                failed: "Gick fel"
            },
            label: "Error"
        }
    }
};



let tFactory = (lang) => {
    if(!translations[lang]) {
        // fallback to english if the language doesn't exist
        lang = "en";
    }

    return (key) => {
        return _.get(translations[lang], key) || _.get(translations["en"], key) || (
            console.warn(`Translation key missing even in English: ${key}`),
                key
        );
    };
};

export {
    tFactory,
    formatCurrency,
    makeLocale
}