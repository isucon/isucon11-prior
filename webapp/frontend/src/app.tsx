import { AppBar, Button, Container, CssBaseline, makeStyles, Toolbar, Typography } from '@material-ui/core';
import React from 'react';
import { VFC } from "react";
import { Route, Switch, useHistory } from 'react-router'
import { LockOpen, PersonAdd } from '@material-ui/icons'
import { RootPage } from './rootPage';
import { SignupPage } from './signupPage';
import { LoginPage } from './loginPage';
import { useAppContext } from './context';
import { SchedulePage } from './schedulePage';

const useStyles = makeStyles(() => ({
  title: {
    flexGrow: 1,
  }
}))

export const App: VFC = () => {
  const classes = useStyles();
  const history = useHistory();
  const { user, logout } = useAppContext();

  return (
    <>
      <CssBaseline />
      <AppBar position="static">
        <Container>
          <Toolbar>
            <Typography variant="h6" className={classes.title}>ISUCON 2021 Prior</Typography>

            { user ? (
              <Button color="inherit" onClick={() => logout()}>{user.nickname}</Button>
            ): (
              <>
                <Button color="inherit" onClick={() => history.push('/signup')} startIcon={<PersonAdd />}>Signup</Button>
                <Button color="inherit" onClick={() => history.push('/login')} startIcon={<LockOpen />}>Login</Button>
              </>
            )}
          </Toolbar>
        </Container>
      </AppBar>
      <Switch>
        <Route path="/" exact component={RootPage} />
        <Route path="/signup" exact component={SignupPage} />
        <Route path="/login" exact component={LoginPage} />
        <Route path="/schedules/:id" exact component={SchedulePage} />
      </Switch>
    </>
  );
}
