import { Button, ButtonGroup, Container, createStyles, FormControlLabel, FormGroup, List, ListItem, ListItemText, makeStyles, Paper, Switch, Typography } from "@material-ui/core";
import React, { useEffect, useState, VFC } from "react";
import { useHistory } from "react-router";
import { useAppContext } from "./context";
import { Schedule } from "./model";

const useStyles = makeStyles((theme) => createStyles({
  empty: {
    padding: theme.spacing(4),
  },
  pager: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
  switcherGroup: {
    marginBottom: theme.spacing(2),
  },
  switcher: {
    justifyContent: 'end',
  },
}))

const PageCount = 30;

export const RootPage: VFC = () => {
  const classes = useStyles();
  const history = useHistory();
  const { user } = useAppContext();
  const [state, update] = useState<{schedules: Schedule[], page: number, loading: boolean, all: boolean }>({
    schedules: [],
    page: 0,
    loading: true,
    all: false,
  });

  useEffect(() => {
    update((s) => ({ ...s, loading: true }));
    let url = '/api/schedules';
    if (state.all) {
      url += '?reserved=1'
    }
    fetch(url).then((r) => r.json()).then((res) => {
      update((s) => ({ ...s, schedules: res, loading: false }));
    })
  }, [update, state.all])

  const maxPage = Math.floor(state.schedules.length / PageCount) + 1;
  let { page } = state;
  if (page <= 0) page = 0;
  if (page >= maxPage) page = maxPage - 1;
  const schedules = state.schedules.slice(page * PageCount, page * PageCount + PageCount)
  const hasPrev = page > 0;
  const hasNext = state.schedules.length > page * PageCount + PageCount;

  const Pagination = () => (
    <ButtonGroup fullWidth className={classes.pager} variant="contained">
      <Button onClick={() => update((s) => ({ ...s, page: 0 }))} disabled={state.page <= 0}>First</Button>
      <Button onClick={() => update((s) => ({ ...s, page: s.page - 1 }))} disabled={!hasPrev}>Prev</Button>
      <Button disabled>{page + 1} / {maxPage}</Button>
      <Button onClick={() => update((s) => ({ ...s, page: s.page + 1 }))} disabled={!hasNext}>Next</Button>
      <Button onClick={() => update((s) => ({ ...s, page: maxPage - 1 }))} disabled={state.page >= maxPage - 1}>Last</Button>
    </ButtonGroup>
  )

  return (
    <Container>
      <Pagination />
      { user && user.staff ? (
        <FormGroup className={classes.switcherGroup}>
          <FormControlLabel
            className={classes.switcher}
            control={<Switch size="small" onChange={(e) => {
              update((s) => ({ ...s, all: e.target.checked }));
            }} />}
            label="Show reserved schedules"
          />
        </FormGroup>
      ) : null}
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
            {state.loading ? 'Loading' : 'No schedule'}
          </Typography>
        </Paper>
      )}
      <Pagination />
    </Container>
  )
}
