import React, {useState} from 'react';
import Paper from "@material-ui/core/Paper";
import {makeStyles, Slider, Typography} from "@material-ui/core";
import {tFactory} from "../../util/i18n";
import Grid from "@material-ui/core/Grid";
import Button from "@material-ui/core/Button";
import {topup} from "../../services/billing";

const useStyles = makeStyles(theme => ({
    page: {
        margin: "2rem",
        padding: "2rem"
    }
}));

let TopupPage = ({language, cancelled, success, stripe}) => {
    let t = tFactory(language);

    let [topupAmount, setTopupAmount] = useState(5);

    let classes = useStyles();

    let onTopupClick = () => {
        topup(topupAmount).then((topupResponse) => {
            stripe.redirectToCheckout({
                sessionId: topupResponse.stripe_session_id
            }).then((result) => {

            })
        });
    };

    return (
        <Paper className={classes.page}>
            <Grid container>
                <Grid item xs={12}>
                    <Typography>{t('top-up.title')}</Typography>
                </Grid>
                <Grid item xs={6}>
                    <Typography gutterBottom>
                        {t('top-up.amount-label')}
                    </Typography>
                </Grid>
                <Grid item xs={6}>
                    <Slider
                        defaultValue={10}
                        value={topupAmount}
                        onChange={(e, v) => setTopupAmount(v)}
                        step={null}
                        min={500}
                        max={5000}
                        valueLabelDisplay="off"
                        marks={[
                            {
                                value: 500,
                                label: "£5"
                            },
                            {
                                value: 1000,
                                label: "£10"
                            },
                            {
                                value: 1500,
                                label: "£15"
                            },
                            {
                                value: 2500,
                                label: "£25"
                            },
                            {
                                value: 5000,
                                label: "£50"
                            }
                        ]}
                    />
                </Grid>
                <Grid item xs={6}>
                    <Button
                        variant="contained"
                        color="primary"
                        onClick={onTopupClick}
                    >{t('cta.topup')}</Button>
                </Grid>
            </Grid>
        </Paper>
    )
};

export {
    TopupPage
}

export default TopupPage;