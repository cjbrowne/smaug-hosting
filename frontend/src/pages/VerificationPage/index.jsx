import React, {useEffect, useState} from 'react';
import {CircularProgress, Paper} from "@material-ui/core";
import {Link} from "react-router-dom";
import {verify} from "../../services/auth";
import qs from 'qs';

let VerificationPage = () => {
    let [success, setSuccess] = useState(false);
    let [requestSent, setRequestSent] = useState(false);

    useEffect(() => {
        let search = qs.parse(window.location.search.substr(1));
        verify(search.code).then((result) => {
            setSuccess(true);
        }, (reason) => {
            setSuccess(false);
        }).finally(() => {
            setRequestSent(true);
        })
    }, []);

    return (
        <React.Fragment>
            {requestSent ?
                success ?
                    <Paper>
                        Thank you for verifying your email address!<br/>
                        <Link to="/login">Log in now</Link>
                    </Paper> :
                    <Paper>
                        That verification link didn't work... sorry. Please <a
                        href="mailto:support@smaug-hosting.co.uk">Contact
                        Support</a><br/>
                    </Paper>
                :
                <Paper>
                    Verifying... <CircularProgress/>
                </Paper>
            }
        </React.Fragment>
    );
};

export {
    VerificationPage
}

export default VerificationPage