import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { PasswordRecoveryComponent } from 'src/app/shared/password-recovery/password-recovery.component';
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
        data: { title: 'Cadastrar' }
      },
      {
        path: 'recuperar-senha/:token',
        component: PasswordRecoveryComponent,
        data: { title: 'Recuperar Senha' }
      }
  ]
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class LandingRoutingModule { }
