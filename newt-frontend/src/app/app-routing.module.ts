import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { AuthGuard } from './services/auth/auth.guard';
import { AuthInterceptor } from './services/auth/auth.interceptor';
import { ArtifactsComponent } from './shared/artifacts/artifacts.component';
import { SystemLayoutComponent } from './shared/system-layout/system-layout.component';
import { Views } from './util/views.enum';


const routes: Routes = [
  {
    path: '',
    component: SystemLayoutComponent,
    canActivate: [AuthGuard],
    children: [
      {
        path: Views.artifacts.url,
        component: ArtifactsComponent
      }
    ]
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  providers: [AuthGuard],
  exports: [RouterModule]
})
export class AppRoutingModule { }
