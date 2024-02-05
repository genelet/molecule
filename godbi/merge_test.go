package godbi

/*
import (
	"testing"
)

func TestmergeMap(t *testing.T) {
	hash1 := map[string]interface{}{"key1":"str1","key2":"str2","key3":3,"key4":false}
	hash2 := map[string]interface{}{"key5":"str5","key6":"str6","key7":7,"key8":true}
	mergedMap := mergeMap(hash1,hash2)
	t.Errorf("%#v", cloneMap(hash2))
	t.Errorf("%#v", mergedMap)
}
func TestmergeArgsMap(t *testing.T) {
	args1 := map[string]interface{}{"key1":"str1","key2":"str2","key3":3,"key4":false}
	args2 := map[string]interface{}{"key5":"str5","key6":"str6","key7":7,"key8":true}
	mergedArgs := mergeArgs(args1,args2)
	t.Errorf("%#v", cloneArgs(args2))
	t.Errorf("%#v", mergedArgs)

	args3 := []interface{}{map[string]interface{}{"key11":"str11", "key12":"str12", "key13":13}, map[string]interface{}{"key21":"str21", "key22":"str22", "key23":23}}
	mergedArgs = mergeArgs(args1,args3)
	t.Errorf("%#v", cloneArgs(args3))
	t.Errorf("%#v", mergedArgs)

	args4 := []map[string]interface{}{map[string]interface{}{"key11":"str11", "key12":"str12", "key13":13}, map[string]interface{}{"key21":"str21", "key22":"str22", "key23":23}}
	mergedArgs = mergeArgs(args1,args4)
	t.Errorf("%#v", mergedArgs)
}

func TestmergeArgs(t *testing.T) {
	args3 := []map[string]interface{}{map[string]interface{}{"key31":"str31", "key32":"str32", "key33":33}, map[string]interface{}{"key311":"str311", "key312":"str312", "key313":313}}
	args4 := []map[string]interface{}{map[string]interface{}{"key41":"str41", "key42":"str42", "key43":43}, map[string]interface{}{"key411":"str411", "key412":"str412", "key413":413}}
	args5 := []interface{}{           map[string]interface{}{"key51":"str51", "key52":"str52", "key53":53}, map[string]interface{}{"key511":"str511", "key512":"str512", "key513":513}}
	args6 := []interface{}{           map[string]interface{}{"key61":"str61", "key62":"str62", "key63":63}, map[string]interface{}{"key611":"str611", "key612":"str612", "key613":613}}

	mergedArgs := mergeArgs(args3,args4)
	t.Errorf("%#v", mergedArgs)

	mergedArgs  = mergeArgs(args3,args5)
	t.Errorf("%#v", mergedArgs)

	mergedArgs  = mergeArgs(args4,args5)
	t.Errorf("%#v", mergedArgs)

	mergedArgs  = mergeArgs(args5,args6)
	t.Errorf("%#v", mergedArgs)

	args4  = []map[string]interface{}{map[string]interface{}{"key31":"str31", "key32":"str32", "key33":33}, map[string]interface{}{"key311":"str311", "key312":"str312", "key313":313}}
	args5  = []interface{}{           map[string]interface{}{"key31":"str31", "key32":"str32", "key33":33}, map[string]interface{}{"key311":"str311", "key312":"str312", "key313":313}}

	mergedArgs  = mergeArgs(args3,args4)
	t.Errorf("%#v", mergedArgs)

	mergedArgs  = mergeArgs(args3,args5)
	t.Errorf("%#v", mergedArgs)
}
*/
