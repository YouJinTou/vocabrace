import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { NgModule } from '@angular/core';
import { MatDialogModule } from '@angular/material/dialog';
import { MatSelectModule } from '@angular/material/select';
import { MatDividerModule } from '@angular/material/divider';
import { HttpClientModule } from '@angular/common/http';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { BlanksDialog, ScrabbleComponent } from './components/scrabble/scrabble.component';
import { WebsocketService } from 'src/services/websocket.service';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { GamesOverviewComponent } from './components/games-overview/games-overview.component';
import { ScrabbleOverviewComponent } from './components/games-overview/scrabble-overview/scrabble-overview.component';
import { FacebookComponent } from './components/login/facebook/facebook.component';
import { LoginStatusService } from 'src/services/login-status.service';
import { HeaderComponent } from './components/header/header.component';

@NgModule({
  declarations: [
    AppComponent,
    ScrabbleComponent,
    BlanksDialog,
    GamesOverviewComponent,
    ScrabbleOverviewComponent,
    FacebookComponent,
    HeaderComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    MatDialogModule,
    MatSelectModule,
    MatDividerModule,
    HttpClientModule
  ],
  providers: [WebsocketService, LoginStatusService],
  bootstrap: [AppComponent]
})
export class AppModule { }
