import React, {useEffect, useState} from 'react';
import _ from 'lodash';
import {
    Button,
    Container,
    makeStyles,
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableRow,
    TableSortLabel,
    IconButton
} from "@material-ui/core";
import {tFactory} from "../../util/i18n";
import {destroy, getAll, start, stop} from "../../services/container";
import Error from "../../components/Error";
import {Mood, MoodBad, SentimentSatisfiedAlt, Assignment as Clipboard} from "@material-ui/icons";
import {useInterval} from "../../util/hooks";

const useStyles = makeStyles(theme => ({
    mood: {
        color: "#00ff00",
        marginRight: "0.5rem"
    },
    moodBad: {
        color: "#ff0000",
        marginRight: "0.5rem"
    },
    moodContainer: {
        display: "flex",
        alignItems: "center"
    },
    neutral: {
        color: "#ffff00",
        marginRight: "0.5rem"
    }
}));

let WhelpDashboard = ({language}) => {
    let [orderBy, setOrderBy] = useState('name');
    let [orderDir, setOrderDir] = useState("asc");
    let [unsortedRows, setUnsortedRows] = useState([]);
    let [error, setError] = useState(null);
    let [confirmDelete, setConfirmDelete] = useState(null);
    let [copied, setCopied] = useState(null);


    let classes = useStyles();

    let t = tFactory(language);

    let createSortByHandler = (name) => {
        return () => {
            if(orderBy === name) {
                setOrderDir((orderDir === "asc" ? "desc" : "asc"));
            } else {
                setOrderBy(name);
                setOrderDir("asc");
            }
        }
    };

    let updateRows = () => {
        getAll().then((containers) => {
            setUnsortedRows(containers);
            setError(null);
        }, ({message:reason}) => {
            setError(<Error message={t('errors.fetch-whelps')} details={reason} />);
        })
    };

    useEffect(updateRows, []);

    useInterval(updateRows, 1000);

    let sort = (rows) => {
        return _.orderBy(rows, [orderBy], [orderDir]);
    };

    let deleteWhelp = (id) => {
        if (confirmDelete === id) {
            destroy(id).then(() => {
                updateRows();
                setError(null);
            }, ({message:reason}) => {
                setError(<Error message={t('errors.delete-whelp')} details={reason} />);
                setConfirmDelete(null);
            });
        } else {
            setConfirmDelete(id);
        }
    };

    let stopWhelp = (id) => {
        stop(id).then(() => {
            updateRows();
            setError(null);
        }, ({message:reason}) => {
            setError(<Error message={t('errors.stop-whelp')} details={reason}/>);
        })
    };

    let startWhelp = (id) => {
        start(id).then(() => {
            updateRows();
        }, ({message:reason}) => {
            setError(<Error message={t('errors.start-whelp')} details={reason}/>);
        })
    };

    let copyIpPort = (ip, port, id) => {
        navigator.clipboard.writeText(`${ip}:${port}`).then(() => {
            setCopied(id);
            setTimeout(() => {
                setCopied(null);
            }, 3000);
        });
    };

    return (
        <React.Fragment>
        {error}
        <Table>
            <TableHead>
                <TableRow>
                    <TableCell>
                        <TableSortLabel active={orderBy === 'name'} direction={orderDir} onClick={createSortByHandler('name')}>
                            {t('generic.name')}
                        </TableSortLabel>
                    </TableCell>
                    <TableCell>
                        <TableSortLabel active={orderBy === 'software'} direction={orderDir} onClick={createSortByHandler('software')}>
                            {t('generic.software')}
                        </TableSortLabel>
                    </TableCell>
                    <TableCell>
                        <TableSortLabel active={orderBy === 'tier'} direction={orderDir} onClick={createSortByHandler('tier')}>
                            {t('generic.tier')}
                        </TableSortLabel>
                    </TableCell>
                    <TableCell>
                        {t('generic.status')}
                    </TableCell>
                    <TableCell>
                        {t('generic.ip-port')}
                    </TableCell>
                    <TableCell>
                        {t('generic.actions')}
                    </TableCell>

                </TableRow>
            </TableHead>
            <TableBody>
                {sort(unsortedRows)
                    .map((row, index) => {
                        return (
                            <TableRow key={row.id}>
                                <TableCell>
                                    {row.name}
                                </TableCell>
                                <TableCell>
                                    {_.capitalize(row.software)}
                                </TableCell>
                                <TableCell>
                                    {t(`containers.tiers.${row.tier}.name`)}
                                </TableCell>
                                <TableCell>
                                    <Container className={classes.moodContainer}>
                                    {
                                        row.status.up ?
                                            <Mood className={classes.mood}/>:
                                            row.status.state === "starting" ?
                                                <SentimentSatisfiedAlt className={classes.neutral}/>:
                                                <MoodBad className={classes.moodBad}/>
                                    }
                                        ({row.status.state})
                                    </Container>
                                </TableCell>
                                <TableCell>
                                    {
                                        row.status.up ?

                                        <React.Fragment>
                                                {row.ip}:{row.port}<IconButton onClick={() => copyIpPort(row.ip, row.port, row.id)}><Clipboard /></IconButton>
                                            {
                                                copied === row.id ?
                                                    "Copied!" :
                                                    null
                                            }
                                        </React.Fragment> :
                                        <React.Fragment>
                                            {t('whelps.down')}
                                        </React.Fragment>
                                    }
                                </TableCell>
                                <TableCell>
                                    {!row.status.up &&
                                        <Button variant="contained" color={confirmDelete === row.id ? "secondary" : "primary"} onClick={() => deleteWhelp(row.id)}>
                                        {
                                            confirmDelete === row.id ?
                                            t('cta.confirm') :
                                            t('cta.delete')

                                        }
                                        </Button>
                                    }
                                    <Button variant="contained" color="default" onClick={() => {
                                        if (row.status.up) {
                                            stopWhelp(row.id);
                                        } else {
                                            startWhelp(row.id);
                                        }
                                    }}>
                                        {
                                            row.status.up ?
                                                t('cta.stop') :
                                                t('cta.start')
                                        }
                                    </Button>
                                </TableCell>
                            </TableRow>
                        )
                    })
                }
            </TableBody>
        </Table>
        </React.Fragment>
    )
};

export {
    WhelpDashboard
}

export default WhelpDashboard;