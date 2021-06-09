import { Button, Card, CardContent, Container, createStyles, Grid, makeStyles, TextField, Theme } from "@material-ui/core";
import { Alert } from '@material-ui/lab';
import React, { ChangeEventHandler, FormEventHandler, useCallback, useState, VFC } from "react";
import { useHistory } from "react-router";
import { useAppContext } from "./context";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
    },
    form: {
      marginTop: theme.spacing(5),
    },
  }),
);

export const LoginPage: VFC = () => {
  const history = useHistory();
  const { login } = useAppContext();
  const [error, updateError] = useState('');
  const [form, updateForm] = useState({ email: '' });
  const classes = useStyles();

  const onChangeEmail = useCallback<ChangeEventHandler<HTMLInputElement>>((e) => {
    updateForm((state) => ({...state, email: e.target.value }))
  }, [updateForm])

  const onSubmit = useCallback<FormEventHandler>((e) => {
    e.preventDefault();

    const fd = new FormData();
    fd.append('email', form.email);

    fetch('/api/login', {
      method: 'POST',
      mode: 'same-origin',
      body: fd
    }).then(async (res) => {
      if(res.status === 200) {
        const user = await res.json();
        login(user);
        history.replace('/');
      } else {
        updateError('Login failed');
      }
    })
  }, [form]);

  return (
    <Container maxWidth="sm">
      <Card className={classes.form}>
        <CardContent>
          <form className={classes.root} onSubmit={onSubmit}>
            <Grid container spacing={2}>
              {error ? (
                <Grid item xs={12}>
                  <Alert severity="error">{error}</Alert>
                </Grid>
              ) : null}
              <Grid item xs={12}>
                <TextField label="E-Mail" type="email" fullWidth required onChange={onChangeEmail} variant="outlined" />
              </Grid>
              <Grid item xs={12}>
                <Button color="primary" variant="contained" type="submit" fullWidth>Login</Button>
              </Grid>
            </Grid>
          </form>
        </CardContent>
      </Card>
    </Container>
  )
}
