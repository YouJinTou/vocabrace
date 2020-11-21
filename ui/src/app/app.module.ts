import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { CellComponent } from './components/scrabble/cell/cell.component';
import { TileComponent } from './components/scrabble/tile/tile.component';
import { BoardComponent } from './components/scrabble/board/board.component';
import { GameComponent } from './components/scrabble/game/game/game.component';
import { WebsocketService } from 'src/services/websocket.service';

@NgModule({
  declarations: [
    AppComponent,
    CellComponent,
    TileComponent,
    BoardComponent,
    GameComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule
  ],
  providers: [WebsocketService],
  bootstrap: [AppComponent]
})
export class AppModule { }
