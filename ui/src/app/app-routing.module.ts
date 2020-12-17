import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { GamesOverviewComponent } from './components/games-overview/games-overview.component';
import { WordlinesComponent } from './components/wordlines/wordlines.component';

const routes: Routes = [
  { path: '', component: GamesOverviewComponent },
  { path: 'wordlines/:poolId', component: WordlinesComponent },
  { path: '**', component: GamesOverviewComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
