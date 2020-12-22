import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { CanDeactivateGuard } from 'src/services/can-deactivate-guard.service';
import { GamesOverviewComponent } from './components/games-overview/games-overview.component';
import { WordlinesComponent } from './components/wordlines/wordlines.component';

const routes: Routes = [
  { path: '', component: GamesOverviewComponent },
  { path: 'wordlines/:poolId', component: WordlinesComponent, canDeactivate: [CanDeactivateGuard] },
  { path: '**', component: GamesOverviewComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
