import { Container, createStyles, List, ListItem, ListItemText, makeStyles, Paper, Typography } from "@material-ui/core";
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

export const RootPage: VFC = () => {
  const classes = useStyles();
  const history = useHistory();
  const [schedules, updateSchedules] = useState<Schedule[]>([]);

  useEffect(() => {
    fetch('/api/schedules').then((r) => r.json()).then((res) => {
      updateSchedules(res);
    })
  }, [updateSchedules])

  return (
    <Container className={classes.container}>
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
          <Typography align="center">現在実施予定はありません</Typography>
        </Paper>
      )}
    </Container>
  )
}
