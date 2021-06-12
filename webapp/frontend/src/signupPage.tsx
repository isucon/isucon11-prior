import { Button, Card, CardContent, Container, createStyles, Grid, makeStyles, TextField, Theme } from "@material-ui/core";
import React, { ChangeEventHandler, FormEventHandler, useCallback, useState, VFC } from "react";
import { useHistory } from "react-router";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
    },
    form: {
      marginTop: theme.spacing(5),
    },
  }),
);

export const SignupPage: VFC = () => {
  const history = useHistory();
  const [form, updateForm] = useState({ email: '', nickname: '' });
  const classes = useStyles();

  const onChangeEmail = useCallback<ChangeEventHandler<HTMLInputElement>>((e) => {
    updateForm((state) => ({...state, email: e.target.value }))
  }, [updateForm])
  const onChangeNickname = useCallback<ChangeEventHandler<HTMLInputElement>>((e) => {
    updateForm((state) => ({...state, nickname: e.target.value }))
  }, [updateForm])

  const onSubmit = useCallback<FormEventHandler>((e) => {
    e.preventDefault();

    const fd = new FormData();
    fd.append('email', form.email);
    fd.append('nickname', form.nickname);

    fetch('/api/signup', {
      method: 'POST',
      mode: 'same-origin',
      body: fd
    }).then((res) => {
      if(res.status === 200) {
        history.replace('/');
      }
    })
  }, [form]);

  return (
    <Container maxWidth="sm">
      <Card className={classes.form}>
        <CardContent>
          <form className={classes.root} onSubmit={onSubmit}>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <TextField label="E-Mail" type="email" fullWidth required onChange={onChangeEmail} variant="outlined" />
              </Grid>
              <Grid item xs={12}>
                <TextField label="Nickname" type="nickname" fullWidth required onChange={onChangeNickname} variant="outlined" />
              </Grid>
              <Grid item xs={12}>
                <Button color="primary" variant="contained" type="submit" fullWidth>Sign up</Button>
              </Grid>
            </Grid>
          </form>
        </CardContent>
      </Card>
    </Container>
  )
}
