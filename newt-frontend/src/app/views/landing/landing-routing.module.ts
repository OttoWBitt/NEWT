import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { SignUpComponent } from 'src/app/shared/sign-up/sign-up.component';
import { LandingComponent } from './landing.component';


export const routes: Routes = [
  {
    path: '',
    redirectTo: '/login',
    pathMatch: 'full'
  },
  {
    path: '',
    children: [
      {
        path: 'login',
        component: LandingComponent,
        data: { title: 'Login' }
      },
      {
        path: 'cadastrar',
        component: SignUpComponent,
        data: { title: 'Login' }
      }
  ]
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class LandingRoutingModule { }
