import React, {useEffect, useState} from 'react';
import {BrowserRouter as Router, Redirect, Route, Switch} from 'react-router-dom';
import {hasAuth} from './services/auth';
import {Dashboard, Login, VerificationPage} from './pages/';
import {Box, CssBaseline, MenuItem, Select, Typography} from '@material-ui/core';
import {makeStyles, ThemeProvider} from '@material-ui/styles';
import defaultTheme from './theme';
import './App.css';
import {tFactory} from './util/i18n';
import {Get} from "./services/cache";

const useStyles = makeStyles(theme => ({
  langSelect: {
    position: "fixed",
    left: 0,
    top: 0,
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between"
  },
  langSelectLabel: {
    padding: "0.6rem"
  }
}));

function App() {
  let loginRedir = null;
  let classes = useStyles();

  let [language, setLanguage] = useState(Get('language'));
  let t = tFactory(language);

  if (!hasAuth()
      && window.location.href.indexOf("/login") === -1
      && window.location.href.indexOf("/verify") === -1
  ) {
    loginRedir = <Redirect to="/login" />;
  }

  let onSettingsChange = (newSettings) => {
    setLanguage(newSettings.language);
  };

  let stripe = window.Stripe(process.env.REACT_APP_STRIPE_PUBLIC_KEY);

  useEffect(() => {
    onSettingsChange({
      language: Get('language')
    });

  }, []);

  return (
    <ThemeProvider theme={defaultTheme}>
      <CssBaseline />
      <Router>
        <Switch>
          <Route path="/verify" render={() => <VerificationPage />} />
          <Route path="/login" render={() => (
            <React.Fragment>
              <Box className={classes.langSelect}>
                <Typography className={classes.langSelectLabel}>
                  {t('i18n.select-language.label')}
                </Typography>
                <Select
                  value={language}
                  onChange={(e) => setLanguage(e.target.value)}
                >
                  <MenuItem value={'en'}>(British) English</MenuItem>
                  <MenuItem value={'se'}>Swedish</MenuItem>
                </Select>
              </Box>
              <Login language={language} />
            </React.Fragment>

          )} />
          <Route path="/" render={({history}) => <Dashboard stripe={stripe} language={language} history={history} onSettingsChange={onSettingsChange}/>} />
        </Switch>
        {loginRedir}
      </Router>
    </ThemeProvider>
  );
}

export default App;
