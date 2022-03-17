import React from 'react';
import {Box, ExpansionPanel, ExpansionPanelDetails, ExpansionPanelSummary} from '@material-ui/core';

import {ExpandMore} from '@material-ui/icons';
import {tFactory} from '../util/i18n';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
    error: {
        backgroundColor: theme.error
    }
}));

let Error = ({message, details, language}) => {
    let t = tFactory(language);

    let classes = useStyles();

    return (
        <ExpansionPanel
            className={classes.error}
        >
            <ExpansionPanelSummary
                expandIcon={<ExpandMore/>}
            >
                <Box
                >{t('errors.label')}:&nbsp;</Box>
                <Box
                >{message}</Box>
            </ExpansionPanelSummary>
            <ExpansionPanelDetails>
                {details}
            </ExpansionPanelDetails>
        </ExpansionPanel>
    )
};

export {
    Error
}

export default Error;