import { Button, Container, createStyles, Grid, List, ListItem, ListItemText, makeStyles, Paper, Typography } from "@material-ui/core";
import React, { useEffect, useState, VFC } from "react";
import { useHistory } from "react-router";
import { Schedule } from "./model";

const useStyles = makeStyles((theme) => createStyles({
  container: {
    marginTop: theme.spacing(2),
  },
  empty: {
    padding: theme.spacing(4),
  },
}))

const PageCount = 30;

export const RootPage: VFC = () => {
  const classes = useStyles();
  const history = useHistory();
  const [state, update] = useState<{schedules: Schedule[], page: number, loading: boolean }>({
    schedules: [],
    page: 0,
    loading: false,
  });

  useEffect(() => {
    update((s) => ({ ...s, loading: true }));
    fetch('/api/schedules').then((r) => r.json()).then((res) => {
      update((s) => ({ ...s, schedules: res, loading: false }));
    })
  }, [update])

  const schedules = state.schedules.slice(state.page * PageCount, state.page * PageCount + PageCount)
  const hasPrev = state.page > 0;
  const hasNext = state.schedules.length > state.page * PageCount + PageCount;

  const Pagination = () => (
      <Grid container spacing={2}>
        <Grid item xs={6}>
          {hasPrev && <Button fullWidth variant="outlined" onClick={() => update((s) => ({ ...s, page: s.page - 1 }))}>Prev</Button>}
        </Grid>
        <Grid item xs={6}>
          {hasNext && <Button fullWidth variant="outlined" onClick={() => update((s) => ({ ...s, page: s.page + 1 }))}>Next</Button>}
        </Grid>
      </Grid>
  )

  return (
    <Container className={classes.container}>
      <Pagination />
      {schedules.length > 0 ? (
        <Paper>
          <List>
            {schedules.map((schedule) => {
              return (
                <ListItem key={schedule.id} button onClick={(e) => {
                  e.preventDefault();
                  history.push(`/schedules/${schedule.id}`)
                }}>
                  <ListItemText primary={schedule.title} secondary={`${schedule.reserved} / ${schedule.capacity}`} />
                </ListItem>
              );
            })}
          </List>
        </Paper>
      ) : (
        <Paper className={classes.empty}>
          <Typography align="center">
            {state.loading ? 'ロード中' : '現在実施予定はありません'}
          </Typography>
        </Paper>
      )}
      <Pagination />
    </Container>
  )
}
