---
title: 'Overview'
date: 2018-11-28T15:14:39+10:00
weight: 1
---

`pomo` is a simple CLI for using the [Pomodoro Technique](https://en.wikipedia.org/wiki/Pomodoro_Technique). There are [some](https://taskwarrior.org/) [amazing](https://todoist.com/) task management systems but Pomo is more of a task *execution* or [timeboxing](https://en.wikipedia.org/wiki/Timeboxing) system. Pomo helps you track what you did, how long it took you to do it, and how much effort you expect it to take.

## Background

### Pomodoro Technique

The Pomodoro Technique is simple and effective:

  * Decide on a task you want to accomplish
  * Break the task into timed intervals (pomodoros), [approx. 25 min]
  * After each pomodoro take a short break [approx. 3 - 5 min]
  * Once all pomodoros are completed take a longer break [approx 15 - 20 min]
  * Repeat

### Concepts

The Pomo CLI provides an interface for working with three different types of resources: 

* Pomodoros
* Tasks
* Projects

###### Pomodoros

A Pomodoro represents a single interval of energy you put towards working on a goal, it is 
an idealized stint of time in which a person is expected to work without interruption. Although it should
be avoided, it is often not possible to work without distraction so the Pomodoro may be suspended and resumed.
Once the amount of time allocated for the Pomodoro has elapsed a configurable event is sent to notify the
user that they should toggle the Pomo UI and initiate a break.

Both time spent running and paused is recorded in an object such as below:

```json
{
  "id": 0,
  "task_id": 0,
  "start": "2019-08-30T14:49:06.528547715-05:00",
  "run_time": "30m0.000045359s",
  "pause_time": "2m0.000022721s" 
}
```

###### Tasks

Tasks represent a particular goal you want to accomplish, in a programming context this might be the development
of a feature or a bug fix. When you configure a project you allocate the number of blocks of time (pomodoros) 
you expect the task multiplied by a duration. You should determine what duration works best for you by considering
the complexity or even level of interest you have in a given task. When working on something that requires a lot
of concentration you might considering a longer duration like `50m` while a shorter duration like `15m` might
be appropriate for something else. It's best to avoid durations much beyond one hour.


Often times time estimates are difficult to make and it's perfectly appropriate to extend or reduce the time
allocated to a particular task by a few Pomodoros.

```json
{
  "id": 0,
  "project_id": 0,
  "message": "Refactor DBO",
  "tags": {
    "kind": "side-project",
    "urgent": false
  },
  "pomodoros": [
    {
      "id": 0,
      "task_id": 0,
      "start": "2019-08-30T15:06:52.278583486-05:00",
      "run_time": "30m0.000045359s",
      "pause_time": "2m0.000022721s"
    }
  ],
  "duration": "30m0.000000000s"
}
```

###### Projects

Projects allow for the logical organization of tasks into a hierarchical tree of groups and sub-groups. They can
be useful for limiting the terminal output to a particular category of tasks e.g. `side-projects` vs `consulting`
or to break large projects into multiple categories and visualize the scope of work. In pomo `>= 0.8.0` all tasks
are associated with a root project.


![Projects](/pomo.svg)

The order of projects is never fixed so you are free to adjust their organization as needed.

