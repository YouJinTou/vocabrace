import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { NgModule } from '@angular/core';
import { MatDialogModule } from '@angular/material/dialog';
import { MatSelectModule } from '@angular/material/select';
import { MatDividerModule } from '@angular/material/divider';
import { HttpClientModule } from '@angular/common/http';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { WordlinesComponent } from './components/wordlines/wordlines.component';
import { WebsocketService } from 'src/services/websocket.service';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { GamesOverviewComponent } from './components/games-overview/games-overview.component';
import { WordlinesOverviewComponent } from './components/games-overview/wordlines-overview/wordlines-overview.component';
import { FacebookComponent } from './components/external-login/facebook/facebook.component';
import { ContextService } from 'src/services/context.service';
import { HeaderComponent } from './components/header/header.component';
import { ExternalLoginComponent } from './components/external-login/external-login.component';
import { BlanksDialog } from './components/wordlines/blanks/blanks.component';
import { GameOverDialog } from './components/wordlines/game-over/game-over.component';
import { TimerComponent } from './components/timer/timer.component';
import { CookieService } from 'ngx-cookie-service';
import { CookiesComponent } from './components/cookies/cookies.component';

@NgModule({
  declarations: [
    AppComponent,
    WordlinesComponent,
    BlanksDialog,
    GamesOverviewComponent,
    WordlinesOverviewComponent,
    FacebookComponent,
    HeaderComponent,
    ExternalLoginComponent,
    GameOverDialog,
    TimerComponent,
    CookiesComponent,
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
  providers: [WebsocketService, ContextService, CookieService],
  bootstrap: [AppComponent]
})
export class AppModule { }
