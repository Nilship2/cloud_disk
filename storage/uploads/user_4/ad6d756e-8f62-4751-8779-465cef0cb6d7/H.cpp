#include<bits/stdc++.h>
using namespace std;
int a[200005],b[200005];
int main(){
	int t;
	cin>>t;
	while(t--){
		int n,fa=0,fb=0,aa=0,bb=0;
		cin>>n;
		for(int i=1;i<=n;i++)cin>>a[i];
		for(int i=1;i<=n;i++)cin>>b[i];
		for(int i=1;i<=n;i++){
			if(a[i]>b[i]){
				aa+=a[i];
			}
			if(a[i]<b[i]){
				bb+=b[i];
			}
			if(a[i]==b[i]){
			    if(a[i]==-1)fa++;
			    if(a[i]==1)fb++;
			}
		}
		while(fa--){
			if(aa>=bb){
				aa--;
			}else{
				bb--;
			}
		}
		while(fb--){
			if(aa<=bb){
				aa++;
			}else{
				bb++;
			}
		}
		cout<<min(aa,bb)<<endl;
	}
	return 0;
}
