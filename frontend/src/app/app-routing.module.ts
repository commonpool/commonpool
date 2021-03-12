import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {ResourceListViewComponent} from './resources/resource-list-view/resource-list-view.component';
import {CreateOrEditResourceComponent} from './resources/create-or-edit-resource/create-or-edit-resource.component';
import {ResourceDetailsComponent} from './resources/resource-details/resource-details.component';
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
import {HomePageComponent} from './home/home-page/home-page.component';
import {BlocksComponent} from './chat/blocks/blocks.component';
import {SampleComponent} from './sample/sample/sample.component';
import {TradingHistoryComponent} from './trading/history/trading-history.component';

const routes: Routes = [
  {
    path: '',
    component: HomePageComponent
  },
  {
    path: 'resources/search',
    component: ResourceListViewComponent
  }, {
    path: 'resources/new',
    component: CreateOrEditResourceComponent
  },
  {
    path: 'resources/:id',
    component: ResourceDetailsComponent
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
      {path: '', redirectTo: 'posts', pathMatch: 'full'},
      {path: 'groups', component: UserGroupsViewComponent},
      {path: 'posts', component: UserResourcesViewComponent, data: {accountType: 'user'}},
      {path: 'transactions', component: OfferListComponent},
      {path: 'transactions/:id', component: OfferDetailsComponent},
    ]
  }, {
    path: 'messages',
    component: ConversationThreadListComponent,
    children: [
      {
        path: 'c/:id',
        component: ConversationThreadComponent,
        data: {
          type: 'channel'
        }
      },
      {
        path: 'g/:id',
        component: ConversationThreadComponent,
        data: {
          type: 'group'
        }
      }
    ]
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
      {path: '', redirectTo: 'posts', pathMatch: 'full'},
      {path: 'members', component: GroupMembersViewComponent},
      {path: 'invitations', component: GroupInvitesViewComponent},
      {path: 'resources', component: GroupResourcesViewComponent},
      {path: 'posts', component: UserResourcesViewComponent, data: {accountType: 'group'}},
      {
        path: 'needs/:resourceId',
        component: ResourceDetailsComponent,
        data: {accountType: 'group', resourceType: 'needs'}
      },
      {
        path: 'offers/:resourceId',
        component: ResourceDetailsComponent,
        data: {accountType: 'group', resourceType: 'offers'}
      }
    ]
  }, {
    path: 'trading-history',
    component: TradingHistoryComponent
  },
  {
    path: 'sample',
    component: SampleComponent
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
