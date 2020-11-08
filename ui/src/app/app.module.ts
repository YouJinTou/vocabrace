import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { CellComponent } from './components/scrabble/cell/cell.component';
import { TileComponent } from './components/scrabble/tile/tile.component';
import { BoardComponent } from './components/scrabble/board/board.component';

@NgModule({
  declarations: [
    AppComponent,
    CellComponent,
    TileComponent,
    BoardComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
