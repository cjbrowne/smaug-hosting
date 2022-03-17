import React, {useState} from 'react';
import {Button, Container, Grid, MenuItem, Paper, Select, Typography} from "@material-ui/core";
import {tFactory} from "../../util/i18n";
import {Put} from "../../services/cache";

// const useStyles = makeStyles(theme => ({
// }));

let SettingsPage = ({language, onSettingsChange}) => {
    let [languageSetting, setLanguageSetting] = useState(language);

    // let classes = useStyles();

    let updateSetting = (setting) => {
        return (evt) => {
            switch(setting) {
                case 'language':
                    setLanguageSetting(evt.target.value);
                    break;
                default:
                    break;
            }
        }
    };

    let saveSettings = () => {
        if(onSettingsChange && typeof onSettingsChange === "function") {
            onSettingsChange({
                language: languageSetting
            });
        }
        Put('language', languageSetting);
    };

    let t = tFactory(language);

    return (
        <React.Fragment>
            <Container>
                <Paper>
                    <Grid container>
                        <Grid item xs={12}>
                            <Typography variant="h3">{t('settings.title')}</Typography>
                        </Grid>
                        <Grid item xs={6}>
                            <Select
                                value={languageSetting}
                                onChange={updateSetting('language')}
                            >
                                <MenuItem value={'en'}>(British) English</MenuItem>
                                <MenuItem value={'se'}>Swedish</MenuItem>
                                <MenuItem value={'gamer'}>Secret Option (shh!)</MenuItem>
                            </Select>
                        </Grid>
                        <Grid item xs={6}>
                            <Button onClick={saveSettings}>{t('settings.save')}</Button>
                        </Grid>
                    </Grid>
                </Paper>
            </Container>
        </React.Fragment>
    )
};

export {
    SettingsPage
}

export default SettingsPage;