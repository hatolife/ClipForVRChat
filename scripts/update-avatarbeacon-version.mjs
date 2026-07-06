#!/usr/bin/env node
import fs from 'node:fs'
import path from 'node:path'

const root = process.cwd()
const version = process.argv[2]
const checkOnly = process.argv.includes('--check')

if (!version || !/^v?\d+\.\d+\.\d+(-(a|b|rc)\d+)?$|^ci$|^dev$/.test(version)) {
  console.error('Usage: node scripts/update-avatarbeacon-version.mjs <vX.Y.Z[-aN|-bN|-rcN]|ci|dev> [--check]')
  process.exit(2)
}

const avatarRoot = path.join(root, 'avatar-gimmicks', 'AvatarBeacon', 'Assets', 'PoppoWorks', 'AvatarBeacon')
const versionTxtPath = path.join(avatarRoot, 'Version.txt')
const versionMetaPath = `${versionTxtPath}.meta`
const prefabPath = path.join(avatarRoot, 'Prefabs', 'AvatarBeacon.prefab')

const gameObjectID = '910239000000000001'
const transformID = '910239000000000002'
const rootTransformID = '2902941087723332167'
const displayName = `AvatarBeacon Version ${version}`

const versionTxt = [
  `AvatarBeacon Version: ${version}`,
  'This file name is intentionally stable so Unity package imports overwrite it.',
  '',
].join('\n')

const versionMeta = [
  'fileFormatVersion: 2',
  'guid: 8e4b95b8e18f4d5ca9b8d770e8218d7b',
  'TextScriptImporter:',
  '  externalObjects: {}',
  '  userData: ',
  '  assetBundleName: ',
  '  assetBundleVariant: ',
  '',
].join('\n')

function readText(file) {
  return fs.readFileSync(file, 'utf8').replace(/\r\n/g, '\n')
}

function assertOrWrite(file, expected) {
  if (checkOnly) {
    const actual = readText(file)
    if (actual !== expected) {
      console.error(`${file} is not updated for AvatarBeacon ${version}`)
      process.exitCode = 1
    }
    return
  }
  fs.writeFileSync(file, expected, 'utf8')
}

function versionPrefabBlock() {
  return [
    `--- !u!1 &${gameObjectID}`,
    'GameObject:',
    '  m_ObjectHideFlags: 0',
    '  m_CorrespondingSourceObject: {fileID: 0}',
    '  m_PrefabInstance: {fileID: 0}',
    '  m_PrefabAsset: {fileID: 0}',
    '  serializedVersion: 6',
    '  m_Component:',
    `  - component: {fileID: ${transformID}}`,
    '  m_Layer: 0',
    `  m_Name: ${displayName}`,
    '  m_TagString: Untagged',
    '  m_Icon: {fileID: 0}',
    '  m_NavMeshLayer: 0',
    '  m_StaticEditorFlags: 0',
    '  m_IsActive: 0',
    `--- !u!4 &${transformID}`,
    'Transform:',
    '  m_ObjectHideFlags: 0',
    '  m_CorrespondingSourceObject: {fileID: 0}',
    '  m_PrefabInstance: {fileID: 0}',
    '  m_PrefabAsset: {fileID: 0}',
    `  m_GameObject: {fileID: ${gameObjectID}}`,
    '  serializedVersion: 2',
    '  m_LocalRotation: {x: 0, y: 0, z: 0, w: 1}',
    '  m_LocalPosition: {x: 0, y: 0, z: 0}',
    '  m_LocalScale: {x: 1, y: 1, z: 1}',
    '  m_ConstrainProportionsScale: 0',
    '  m_Children: []',
    `  m_Father: {fileID: ${rootTransformID}}`,
    '  m_LocalEulerAnglesHint: {x: 0, y: 0, z: 0}',
    '',
  ].join('\n')
}

function updatePrefab(prefab) {
  const blockPattern = new RegExp(`--- !u!1 &${gameObjectID}\\n[\\s\\S]*?--- !u!4 &${transformID}\\n[\\s\\S]*?(?=\\n--- !u!|\\s*$)`)
  let next = prefab.replace(blockPattern, '')
  const childLine = `  - {fileID: ${transformID}}`
  if (!next.includes(childLine)) {
    const rootPattern = new RegExp(`(--- !u!4 &${rootTransformID}\\n[\\s\\S]*?  m_Children:\\n)((?:  - \\{fileID: \\d+\\}\\n)+)`, 'm')
    next = next.replace(rootPattern, (_match, head, children) => `${head}${children}${childLine}\n`)
  }
  next = next.replace(/\n*$/, '\n')
  return `${next}${versionPrefabBlock()}`
}

assertOrWrite(versionTxtPath, versionTxt)
assertOrWrite(versionMetaPath, versionMeta)

const expectedPrefab = updatePrefab(readText(prefabPath))
assertOrWrite(prefabPath, expectedPrefab)

if (process.exitCode) {
  process.exit(process.exitCode)
}

console.log(`AvatarBeacon version OK: ${version}`)
