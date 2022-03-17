import React, {useState} from 'react';
import {
    Button, CircularProgress,
    Dialog,
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
    Grid, makeStyles,
    MenuItem,
    Select,
    TextField
} from "@material-ui/core";
import {tFactory} from "../util/i18n";
import {green} from "@material-ui/core/colors";

const useStyles = makeStyles(theme => ({
    createButtonWrapper: {
        position: "relative",
    },
    createButtonProgress: {
        color: green[500],
        position: "absolute",
        top: "50%",
        left: "50%",
        marginTop: -12,
        marginLeft: -12
    }
}));

let CreateContainerDialog = ({open, software, language, onClose, onCreateClick}) => {
    let t = tFactory(language);

    let [name, setName] = useState('');
    let [tier, setTier] = useState(0);
    let [creating, setCreating] = useState(false);

    let classes = useStyles();

    let handleCreateClick = () => {
        if (typeof onCreateClick === "function") {
            let res = onCreateClick({name, tier, software});
            if (res instanceof Promise) {
                setCreating(true);
                res.then(() => {
                    setCreating(false);
                });
            } else {
                setCreating(false);
            }
        } else {
            setCreating(false);
        }
    };

    let handleClose = () => {
        // don't close dialog while creation is in process
        if (creating) return;
        onClose();
    };

    return (
        <Dialog open={open} onClose={handleClose}>
            <DialogTitle>{t('containers.create-dialog.title')} - {t(`containers.${software}.title`)}</DialogTitle>
            <DialogContent>
                <Grid container spacing={3}>
                    <Grid item xs={12}>
                        <TextField
                            id="name"
                            label={t('containers.create-dialog.form.name.label')}
                            placeholder={t('containers.create-dialog.form.name.placeholder')}
                            type="text"
                            value={name}
                            onChange={e => setName(e.target.value)}
                        />
                    </Grid>
                    <Grid item xs={4}>
                        <Select
                            label={t('containers.create-dialog.form.tier.label')}
                            value={tier}
                            onChange={e => setTier(e.target.value)}
                        >
                            <MenuItem value={0}>
                                {t('containers.tiers.0.name')}
                            </MenuItem>
                            <MenuItem value={1}>
                                {t('containers.tiers.1.name')}
                            </MenuItem>
                            <MenuItem value={2}>
                                {t('containers.tiers.2.name')}
                            </MenuItem>
                            <MenuItem value={3}>
                                {t('containers.tiers.3.name')}
                            </MenuItem>
                        </Select>
                    </Grid>
                    <Grid item xs>
                        {t(`containers.tiers.${tier}.help`)}
                    </Grid>

                </Grid>
                <DialogContentText>

                </DialogContentText>
                <DialogActions>
                    <div className={classes.createButtonWrapper}>
                        <Button onClick={handleCreateClick} variant="contained" color="primary" disabled={creating}>
                            {t('cta.create')}
                        </Button>
                        {creating && <CircularProgress size={24} className={classes.createButtonProgress}/>}
                    </div>
                    <Button onClick={onClose} variant="contained" color="default" disabled={creating}>
                        {t('cta.cancel')}
                    </Button>
                </DialogActions>
            </DialogContent>
        </Dialog>
    );
};

export default CreateContainerDialog;