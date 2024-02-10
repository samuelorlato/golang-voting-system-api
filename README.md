# Golang Voting System API

A Go API that uses websockets to handle voting

## Routes

#### Vote

<details>
 <summary><code>GET</code> <code><b>/vote</b></code> <code>(creates a room and enter as voter)</code></summary>

##### JSON Body Params

> | name   | type     | data type |
> | ------ | -------- | --------- |
> | option | required | string    |

![/vote](./assets/vote1.png)

</details>

<details>
 <summary><code>GET</code> <code><b>/vote/:roomId</b></code> <code>(enters in a room as voter)</code></summary>

##### JSON Body Params

> | name   | type     | data type |
> | ------ | -------- | --------- |
> | option | required | string    |

![/vote](./assets/vote2.png)

</details>

#### Spectate

<details>
 <summary><code>GET</code> <code><b>/spectate</b></code> <code>(creates a room and enter as spectator)</code></summary>

![/vote](./assets/spectate1.png)

</details>

<details>
 <summary><code>GET</code> <code><b>/spectate/:roomId</b></code> <code>(enters in a room as spectator)</code></summary>

![/vote](./assets/spectate1.png)

</details>
