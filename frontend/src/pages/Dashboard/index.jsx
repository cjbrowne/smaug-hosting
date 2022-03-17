import React, {useEffect, useRef, useState} from 'react';
import {Redirect, Route, Switch} from 'react-router-dom';
import {
    AppBar,
    Box,
    Button,
    Card,
    CardActions,
    CardContent,
    CardHeader,
    CardMedia,
    Container,
    Drawer,
    Grid,
    IconButton,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
    Menu,
    Snackbar,
    SnackbarContent,
    Toolbar,
    Typography
} from '@material-ui/core';
import {Close, Error, ExitToApp, Home, Menu as MenuIcon, Settings, Whatshot} from '@material-ui/icons';
import {formatCurrency, makeLocale, tFactory} from '../../util/i18n';
import {makeStyles} from '@material-ui/styles';
import {logout} from '../../services/auth';
import {create} from '../../services/container';
import {subscribeToBalance} from '../../services/billing';

import CreateContainerDialog from '../../components/CreateContainerDialog';
import WhelpDashboard from './WhelpDashboard';
import SettingsPage from './SettingsPage';
import TopupPage from './TopupPage';
import {green} from "@material-ui/core/colors";
import MenuItem from "@material-ui/core/MenuItem";


let useStyles = makeStyles(theme => ({
    siteTitle: {
        flexGrow: 1
    },
    nav: {
        display: "flex",
        flexDirection: "column",
        height: "100%"
    },
    logout: {
        alignSelf: "flex-end",
        marginTop: "auto"
    },
    cardMedia: {
        height: 0,
        paddingTop: '56.25%', // 16:9
    },
    errorSnackbar: {
        backgroundColor: theme.palette.error.light
    },
    successSnackbar: {
        backgroundColor: green[600]
    },
    snackbarContent: {
        display: 'flex',
        alignItems: 'center'
    },
    emoji: {
        height: "1.5rem",
    }
}));

let Dashboard = ({language = 'en', country = 'gb', history, onSettingsChange, stripe}) => {
    let [isCreateContainerDialogOpen, setCreateContainerDialogOpen] = useState(false);
    let [createContainerSoftware, setCreateContainerSoftware] = useState('');
    let [sideMenuShown, setSideMenuShown] = useState(false);
    let [redirect, setRedirect] = useState(null);
    let [snackbarContent, setSnackbarContent] = useState(null);
    let [snackbarOpen, setSnackbarOpen] = useState(false);
    let [snackbarVariant, setSnackbarVariant] = useState('error');
    let [accountBalance, setAccountBalance] = useState('?.??');
    let [isAccountMenuOpen, setIsAccountMenuOpen] = useState(false);
    let anchorRef = useRef(null);
    let toggleSideMenu = () => setSideMenuShown(!sideMenuShown);

    useEffect(() => {
        let cancelSub = subscribeToBalance((balanceMicroGBP) => {
            let balance = balanceMicroGBP / 1000;
            // if we ever support extra currencies, we need to change this VVV
            let balanceFormatted = `${formatCurrency(makeLocale('en', 'gb'), balance)}`;
            setAccountBalance(balanceFormatted);
        });
        return () => {
            cancelSub();
        }
    }, [country, language]);

    let t = tFactory(language);

    let classes = useStyles();
    let onLogoutClick = () => {
        logout();
        setRedirect(<Redirect to="/login"/>);
    };

    if (redirect) {
        return redirect;
    }

    let onCreateContainerClick = (software) => {
        setCreateContainerDialogOpen(true);
        setCreateContainerSoftware(software);
    };

    let onCreateContainer = ({name, tier, software}) => {
        return create({name, tier, software}).then((result) => {
            setCreateContainerDialogOpen(false);
            setSnackbarOpen(true);
            setSnackbarContent(t('containers.create-success'));
            setSnackbarVariant('success');
        }, (error) => {
            setCreateContainerDialogOpen(false);
            setSnackbarOpen(true);
            setSnackbarContent(error.message);
            setSnackbarVariant('error');
        });
    };

    let navigate = (path) => {
        history.push(path);
        setSideMenuShown(false);
    };

    return (
        <React.Fragment>
            <Drawer open={sideMenuShown} onClose={() => setSideMenuShown(false)}>
                <List component="nav" className={classes.nav}>
                    <ListItem button key="home" onClick={() => navigate('/')}>
                        <ListItemIcon><Home/></ListItemIcon>
                        <ListItemText>{t('nav.home')}</ListItemText>
                    </ListItem>
                    <ListItem button key={"whelps"} onClick={() => navigate('/whelps')}>
                        <ListItemIcon><Whatshot/></ListItemIcon>
                        <ListItemText>{t('nav.whelp-list')}</ListItemText>
                    </ListItem>
                    <ListItem button key={"settings"} onClick={() => navigate('/settings')}>
                        <ListItemIcon><Settings/></ListItemIcon>
                        <ListItemText>{t('nav.settings')}</ListItemText>
                    </ListItem>
                    <ListItem button key="logout" className={classes.logout} onClick={onLogoutClick}>
                        <ListItemIcon><ExitToApp/></ListItemIcon>
                        <ListItemText>{t('nav.logout')}</ListItemText>
                    </ListItem>
                </List>
            </Drawer>
            <AppBar position="static">
                <Toolbar>
                    <IconButton
                        onClick={toggleSideMenu}
                    >
                        <MenuIcon/>
                    </IconButton>
                    <Typography variant="h6" className={classes.siteTitle}>
                        {t('site.title')}
                    </Typography>
                    <div>
                        <Button
                            ref={anchorRef}
                            onClick={() =>
                                setIsAccountMenuOpen(o => !o)
                            }>{t('account.balance')}: {accountBalance}</Button>
                        <Menu
                            open={isAccountMenuOpen}
                            elevation={0}
                            getContentAnchorEl={null}
                            anchorEl={anchorRef.current}
                            onClose={() => setIsAccountMenuOpen(false)}
                            anchorOrigin={{
                                vertical: 'bottom',
                                horizontal: 'center',
                            }}
                            transformOrigin={{
                                vertical: 'top',
                                horizontal: 'center'
                            }}>
                                <MenuItem onClick={() => {
                                    navigate("/topup");
                                    setIsAccountMenuOpen(false);
                                }}>{t('account.topup')}</MenuItem>
                        </Menu>
                    </div>
                </Toolbar>
            </AppBar>
            <Switch>
                <Route path="/topup/cancel" render={() => <TopupPage language={language} stripe={stripe} cancelled />}/>
                <Route path="/topup/success" render={() => <TopupPage language={language} stripe={stripe} success />}/>
                <Route path="/topup" render={() => <TopupPage language={language} stripe={stripe} />}/>
                <Route path="/settings" render={() => <SettingsPage language={language} onSettingsChange={onSettingsChange}/>}/>
                <Route path="/whelps" render={()=><WhelpDashboard language={language}/> }/>
                <Route path="/" exact>
                    <Container>
                        <Box my={3}>
                            <Grid container spacing={3}>
                                <Card>
                                    <CardHeader>
                                        {t('containers.minecraft.title')}
                                    </CardHeader>
                                    <CardMedia
                                        image="/images/minecraft.png"
                                        title={t('containers.minecraft.title')}
                                        className={classes.cardMedia}
                                    />
                                    <CardContent>
                                        {t('containers.minecraft.description')}
                                    </CardContent>
                                    <CardActions>
                                        <Button variant="contained" color="primary"
                                                onClick={() => onCreateContainerClick('minecraft')}>
                                            {t('cta.create-container')}
                                            <Whatshot/>
                                        </Button>
                                    </CardActions>
                                </Card>
                            </Grid>
                        </Box>
                    </Container>
                </Route>
            </Switch>
            <CreateContainerDialog
                open={isCreateContainerDialogOpen}
                software={createContainerSoftware}
                onClose={() => setCreateContainerDialogOpen(false)}
                onCreateClick={onCreateContainer}
            />
            <Snackbar open={snackbarOpen}
                      autoHideDuration={1500}
                      anchorOrigin={{
                          vertical: 'bottom',
                          horizontal: 'left'
                      }}
                      action={[
                          <IconButton key="close" aria-label="close" color="inherit"
                                      onClick={() => setSnackbarOpen(false)}>
                              <Close/>
                          </IconButton>
                      ]}
                      onClose={() => setSnackbarOpen(false)}
            >
                <SnackbarContent
                    className={classes[`${snackbarVariant}Snackbar`]}
                    message={
                        <div className={classes.snackbarContent}>
                            {snackbarVariant === "success" &&
                            <img alt="dab" src="/images/person_dabbing.png" className={classes.emoji}/>}
                            {snackbarVariant === "error" && <Error/>}
                            <span>
                                {snackbarContent}
                            </span>
                        </div>
                    }
                />
            </Snackbar>
        </React.Fragment>
    )
};

export default Dashboard;