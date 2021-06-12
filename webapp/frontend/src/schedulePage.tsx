import { Breadcrumbs, Button, Card, CardContent, Container, createStyles, Grid, Link, List, ListItem, ListItemText, makeStyles, Paper, Typography } from "@material-ui/core";
import React, { MouseEventHandler, useCallback, useEffect, useState, VFC } from "react";
import { useParams } from "react-router";
import { Link as RouterLink } from 'react-router-dom';
import { useAppContext } from "./context";
import { Schedule } from "./model";

const useStyles = makeStyles((theme) =>
  createStyles({
    container: {
      marginTop: theme.spacing(2),
    },
  }),
);

export const SchedulePage: VFC = () => {
  const classes = useStyles();
  const { user } = useAppContext();
  const { id } = useParams<{ id: string; }>()
  const [state, update] = useState<{ schedule?: Schedule, loading: boolean, tick: number }>({
    loading: true,
    tick: 0,
  });

  const onClick = useCallback<MouseEventHandler>((e) => {
    e.preventDefault();

    const fd = new FormData();
    fd.append('schedule_id', id);
    fetch(`/api/reservations`, {
      method: 'POST',
      mode: 'same-origin',
      body: fd,
    }).then(() => {
      update((s) => ({ ...s, tick: s.tick + 1 }));
    })
  }, [id, update])

  useEffect(() => {
    update((s) => ({ ...s, loading: true }))
    fetch(`/api/schedules/${id}`).then((res) => res.json()).then((schedule) => {
      update((state) => ({ ...state, schedule, loading: false }))
    })
  }, [id, update, state.tick])

  const { schedule } = state
  const reservable = schedule && user ? schedule.reservations.reduce((m, r) => m && r.user.id != user.id, true) : false;

  return (
    <Container className={classes.container}>
      {state.loading ? (
        <Paper>
          <Typography display="block" align="center">Loading</Typography>
        </Paper>
      ) : (
        <>
          <Breadcrumbs>
            <Link component={RouterLink} to="/">Top Page</Link>
            {schedule && <Link component={RouterLink} to={`/schedules/${id}`}>{schedule.title}</Link>}
          </Breadcrumbs>

          <Grid container spacing={2} className={classes.container}>
            <Grid item xs={12} md={8}>
              <Card>
                <CardContent>
                  {schedule ? (
                    <>
                      <Typography variant="h5" component="h1">{schedule.title}</Typography>
                      <Typography variant="caption" display="block" gutterBottom>{schedule.reserved} / {schedule.capacity}</Typography>

                      <List>
                        {schedule.reservations.map((res) => (
                          <ListItem key={res.id}>
                            <ListItemText primary={res.user.nickname} secondary={res.user.email} />
                          </ListItem>
                        ))}
                      </List>
                    </>
                  ) : (
                    <Typography variant="h5">Not Found</Typography>
                  )}
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} md={4}>
              <Card>
                <CardContent>
                  {schedule && reservable ? (
                    <Button fullWidth color="primary" variant="contained" size="large" onClick={onClick}>
                      Reserve now
                    </Button>
                  ) : (
                    <Button fullWidth color="primary" variant="contained" disabled size="large">
                      Can't reserve
                    </Button>
                  )}
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        </>
      )}
    </Container>
  )
}
