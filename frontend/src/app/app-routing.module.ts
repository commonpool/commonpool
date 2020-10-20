import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {ResourceListViewComponent} from './resources/resource-list-view/resource-list-view.component';
import {CreateOrEditResourceComponent} from './resources/create-or-edit-resource/create-or-edit-resource.component';
import {ResourceDetailsComponent} from './resources/resource-details/resource-details.component';
import {UserProfileComponent} from './user-profile/user-profile.component';


const routes: Routes = [
  {
    path: '',
    component: ResourceListViewComponent
  }, {
    path: 'resources/new',
    component: CreateOrEditResourceComponent
  }, {
    path: 'resources/:id',
    component: ResourceDetailsComponent
  }, {
    path: 'resources/:id/edit',
    component: CreateOrEditResourceComponent
  }, {
    path: 'profiles/:id',
    component: UserProfileComponent
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
