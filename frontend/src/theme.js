import { createMuiTheme } from '@material-ui/core/styles';

const theme = createMuiTheme({
    palette: {
        type: "dark",
        primary: {
            main: '#993333',
        },
        secondary: {
            main: '#19897b',
        },
        error: {
            main: '#995544',
        },
        background: {
            default: '#333',
        },
    }
});

export default theme;