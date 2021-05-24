import { AppBar, Button, Container, CssBaseline, makeStyles, Toolbar, Typography } from '@material-ui/core';
import React from 'react';
import { VFC } from "react";
import { Switch, useHistory } from 'react-router'

const useStyles = makeStyles(() => ({
  title: {
    flexGrow: 1,
  }
}))

export const App: VFC = () => {
  const classes = useStyles();
  const history = useHistory();

  return (
    <>
      <CssBaseline />
      <AppBar position="static">
        <Container>
          <Toolbar>
            <Typography variant="h6" className={classes.title}>ISUMARK</Typography>
            <Button color="inherit" onClick={() => history.push('/signup')}>Signup</Button>
            <Button color="inherit" onClick={() => history.push('/login')}>Login</Button>
          </Toolbar>
        </Container>
      </AppBar>
      <Switch>
      </Switch>
    </>
  );
}
