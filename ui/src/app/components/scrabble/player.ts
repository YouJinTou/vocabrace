export class Player {
    name: string
    points: number
    yourMove: boolean

    constructor(name: string, points: number, yourMove: boolean) {
        this.name = name;
        this.points = points;
        this.yourMove = yourMove;
    }
}