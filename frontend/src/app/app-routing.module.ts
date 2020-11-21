import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {ResourceListViewComponent} from './resources/resource-list-view/resource-list-view.component';
import {CreateOrEditResourceComponent} from './resources/create-or-edit-resource/create-or-edit-resource.component';
import {ResourceDetailsComponent} from './resources/resource-details/resource-details.component';
import {UserProfileComponent} from './user-profile/user-profile.component';
import {ConversationThreadListComponent} from './chat/conversation-thread-list/conversation-thread-list.component';
import {ResourceInquiryComponent} from './resources/resource-inquiry/resource-inquiry.component';
import {ConversationThreadComponent} from './chat/conversation-thread/conversation-thread.component';
import {OfferListComponent} from './offers/offer-list/offer-list.component';
import {OfferDetailsComponent} from './offers/offer-details/offer-details.component';
import {CreateOfferComponent} from './offers/create-offer/create-offer.component';
import {CreateOrEditGroupComponent} from './groups/create-or-edit-group/create-or-edit-group.component';
import {GroupViewComponent} from './groups/group-view/group-view.component';
import {GroupMembersViewComponent} from './groups/group-members-view/group-members-view.component';
import {UserViewComponent} from './users/user-view/user-view.component';
import {UserGroupsViewComponent} from './users/user-groups-view/user-groups-view.component';
import {UserResourcesViewComponent} from './users/user-resources-view/user-resources-view.component';
import {GroupInvitesViewComponent} from './groups/group-invites-view/group-invites-view.component';
import {GroupResourcesViewComponent} from './groups/group-resources-view/group-resources-view.component';


const routes: Routes = [
  {
    path: '',
    component: ResourceListViewComponent
  }, {
    path: 'resources/new',
    component: CreateOrEditResourceComponent
  }, {
    path: 'resources/:id/inquire',
    component: ResourceInquiryComponent
  }, {
    path: 'resources/:id/edit',
    component: CreateOrEditResourceComponent
  }, {
    path: 'users/:id',
    component: UserViewComponent,
    children: [
      {path: '', redirectTo: 'groups', pathMatch: 'full'},
      {path: 'groups', component: UserGroupsViewComponent},
      {path: 'needs', component: UserResourcesViewComponent},
      {path: 'offers', component: UserResourcesViewComponent},
      {path: 'needs/:resourceId', component: ResourceDetailsComponent},
      {path: 'offers/:resourceId', component: ResourceDetailsComponent}
    ]
  }, {
    path: 'messages',
    component: ConversationThreadListComponent
  }, {
    path: 'messages/:id',
    component: ConversationThreadComponent
  }, {
    path: 'offers',
    component: OfferListComponent
  }, {
    path: 'offers/new',
    component: CreateOfferComponent,
  }, {
    path: 'offers/:id',
    component: OfferDetailsComponent
  }, {
    path: 'groups/new',
    component: CreateOrEditGroupComponent
  }, {
    path: 'groups/:id',
    component: GroupViewComponent,
    children: [
      {path: '', redirectTo: 'members', pathMatch: 'full'},
      {path: 'members', component: GroupMembersViewComponent},
      {path: 'invitations', component: GroupInvitesViewComponent},
      {path: 'resources', component: GroupResourcesViewComponent},
    ]
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
