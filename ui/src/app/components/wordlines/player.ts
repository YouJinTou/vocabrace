export class Player {
    id: string;
    name: string
    points: number

    constructor(id: string, name: string, points: number) {
        this.id = id;
        this.name = name;
        this.points = points;
    }
}