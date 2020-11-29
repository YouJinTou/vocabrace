import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { GamesOverviewComponent } from './components/games-overview/games-overview.component';
import { ScrabbleComponent } from './components/scrabble/scrabble.component';

const routes: Routes = [
  { path: '', component: GamesOverviewComponent },
  { path: 'scrabble/:poolId', component: ScrabbleComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
