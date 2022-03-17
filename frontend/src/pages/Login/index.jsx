import React, { useState } from 'react';
import { Box, TextField, Button, Card, CardContent, CardActions, Typography, makeStyles, CircularProgress } from '@material-ui/core';
import { Redirect } from 'react-router-dom';

import { tFactory } from '../../util/i18n';

import { login, register, hasAuth } from '../../services/auth';

import { Error } from '../../components/Error';

const useStyles = makeStyles(theme => ({
    formContainer: {
        display: 'flex',
        flexDirection: 'column'
    },
    page: {
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        width: '100vw',
        height: '100vh'
    }
}));

let Login = ({language = "en"}) => {
    let [email, setEmail] = useState('');
    let [password, setPassword] = useState('');
    let [result, setResult] = useState(null);
    let [loggingIn, setLoggingIn] = useState(false);
    let [registering, setRegistering] = useState(false);
    let classes = useStyles();

    let t = tFactory(language);

    let onRegisterClick = () => {
        setLoggingIn(false);
        setRegistering(true);
        register(email, password).then(() => {
            setRegistering(false);
            setResult("Success!  You can now log in.");
        }, (error) => {
            setRegistering(false);
            setResult(
                <Error
                    message={t('errors.register.failed')}
                    details={error.message}
                    language={language}
                    />
            );
        })
    };

    let onLoginClick = () => {
        setLoggingIn(true);
        setRegistering(false);
        setResult(null);
        login(email, password).then(() => {
            setLoggingIn(false);
            setResult(<Redirect to="/" />);
        }, (error) => {
            setLoggingIn(false);
            setResult(<Error 
                            message={t('errors.login.failed')} 
                            details={error.message} 
                            language={language}
                            />);
        });
    };

    if(hasAuth()) {
        return <Redirect to="/" />;
    }

    return (
        <Box className={classes.page}>
            <Card>
                <CardContent>
                    <Typography>
                        {t('form.auth.title')}
                    </Typography>
                    {result}
                    <form autoComplete="off" noValidate className={classes.formContainer}>
                        <TextField 
                            id='login-email'
                            label={t('form.auth.email.label')}
                            type="email"
                            value={email}
                            onChange={(evt) => setEmail(evt.target.value)}
                            margin="normal"
                            variant="outlined"
                            autoComplete="username"
                        />
                        <TextField 
                            id="login-password"
                            label={t('form.auth.password.label')}
                            type="password"
                            value={password}
                            onChange={(evt) => setPassword(evt.target.value)}
                            margin="normal"
                            variant="outlined"
                            autoComplete="password"
                        />
                </form>        
                </CardContent>
                <CardActions>
                    <Button
                        variant="contained"
                        color="primary"
                        onClick={onLoginClick}
                        disabled={registering || loggingIn}
                    >
                        {loggingIn ?
                            <CircularProgress color="inherit" /> : 
                            t('cta.login')}
                    </Button>
                    <Button
                        variant="contained"
                        color="secondary"
                        onClick={onRegisterClick}
                        disabled={registering || loggingIn}
                        >
                        {registering ?
                            <CircularProgress color="inherit" /> :
                            t('cta.register')}
                    </Button>
                </CardActions>
                
            </Card>
        </Box>
    )
}

export default Login;